#!/usr/bin/env node
'use strict';
// postinstall: fetch prebuilt Go binary from GitHub Releases (no Go required).
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
    const req = https.get(url, { headers: { 'User-Agent': 'croctui-installer' } }, res => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        fetch(res.headers.location, dest).then(resolve, reject);
        return;
      }
      if (res.statusCode !== 200) { res.resume(); reject(new Error(`${res.statusCode} ${res.statusMessage}`)); return; }
      const chunks = [];
      res.on('data', c => chunks.push(c));
      res.on('error', reject);
      res.on('end', () => {
        fs.promises.mkdir(path.dirname(dest), { recursive: true }).then(() => {
          return fs.promises.writeFile(dest, Buffer.concat(chunks));
        }).then(resolve, reject);
      });
    });
    req.on('error', reject);
    req.setTimeout(30000, () => req.destroy(new Error('timeout fetching ' + url)));
  });
}

async function main() {
  const bin = binPath();
  if (fs.existsSync(bin)) return;
  let version = '0.1.0';
  try { version = require('../package.json').version; } catch {}
  const tag = `v${version}`;
  const asset = binName();
  const base = process.env.CROCTUI_RELEASE_BASE || 'https://github.com/asterxsk/croctui/releases/download';
  const url = `${base}/${tag}/${asset}`;
  const tmp = bin + '.tmp';
  try {
    await fetch(url, tmp);
    try { await fs.promises.unlink(bin); } catch {}
    try { await fs.promises.rename(tmp, bin); } catch { await fs.promises.copyFile(tmp, bin); await fs.promises.unlink(tmp); }
    if (process.platform !== 'win32') await fs.promises.chmod(bin, 0o755);
    console.log(`croctui: downloaded ${asset} ${tag}`);
  } catch (e) {
    try { await fs.promises.unlink(tmp); } catch {}
    if (!process.env.CI) console.warn(`croctui: no prebuilt binary at ${url} (${e.message}) — will build on first run if Go is available.`);
  }
}
main();
