#!/usr/bin/env node
'use strict';
// postinstall: fetch prebuilt Go binary from GitHub Releases (no Go required).
// Skips download if binary already present or source checkout with Go can build.
const fs = require('node:fs');
const path = require('node:path');
const https = require('node:https');

const PKG_ROOT = path.resolve(__dirname, '..');

function binName() {
  const plat = process.platform === 'win32' ? 'windows' : process.platform;
  const arch = process.arch === 'x64' ? 'amd64' : process.arch === 'arm64' ? 'arm64' : process.arch;
  const ext = process.platform === 'win32' ? '.exe' : '';
  return `croctui-${plat}-${arch}${ext}`;
}

function binPath() { return path.join(__dirname, binName()); }

function fetch(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, { headers: { 'User-Agent': 'croctui-installer' } }, res => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        file.close(() => fs.unlink(dest, () => {}));
        fetch(res.headers.location, dest).then(resolve, reject);
        return;
      }
      if (res.statusCode !== 200) {
        file.close(() => fs.unlink(dest, () => {}));
        reject(new Error(`${res.statusCode} ${res.statusMessage}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => file.close(resolve));
      file.on('error', reject);
    }).on('error', reject);
  });
}

async function main() {
  const bin = binPath();
  if (fs.existsSync(bin)) return; // already present (dev build, cache)

  let version = '0.1.0';
  try { version = require('../package.json').version; } catch {}
  const tag = `v${version}`;
  const asset = binName();
  const base = process.env.CROCTUI_RELEASE_BASE || 'https://github.com/asterxsk/croctui/releases/download';
  const url = `${base}/${tag}/${asset}`;
  const tmp = bin + '.tmp';

  try {
    await fetch(url, tmp);
    fs.renameSync(tmp, bin);
    if (process.platform !== 'win32') fs.chmodSync(bin, 0o755);
    console.log(`croctui: downloaded ${asset} ${tag}`);
  } catch (e) {
    try { fs.unlinkSync(tmp); } catch {}
    // No release yet — silent in dev. Runtime will build from source if Go exists.
    const isCI = !!process.env.CI;
    if (!isCI) console.warn(`croctui: no prebuilt binary at ${url} (${e.message}) — will build on first run if Go is available.`);
  }
}

main();
