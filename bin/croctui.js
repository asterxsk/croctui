#!/usr/bin/env node
'use strict';
const { spawnSync } = require('node:child_process');
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

function localBin() {
  return path.join(__dirname, binName());
}

// Also check legacy names for dev builds.
function fallbackBins() {
  return [
    path.join(__dirname, process.platform === 'win32' ? 'croctui.exe' : 'croctui'),
    path.join(__dirname, process.platform === 'win32' ? 'crocgui.exe' : 'crocgui'),
    path.join(PKG_ROOT, 'bin', process.platform === 'win32' ? 'croctui.exe' : 'croctui'),
  ];
}

function findBin() {
  const primary = localBin();
  if (fs.existsSync(primary)) return primary;
  for (const p of fallbackBins()) if (fs.existsSync(p)) return p;
  return primary;
}

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
        reject(new Error(`download failed: ${res.statusCode} ${res.statusMessage} for ${url}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => file.close(resolve));
      file.on('error', reject);
    }).on('error', reject);
  });
}

async function ensureBin() {
  const bin = findBin();
  if (fs.existsSync(bin)) return bin;

  // Try to download from GitHub Releases (no Go required).
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
    return bin;
  } catch (e) {
    try { fs.unlinkSync(tmp); } catch {}
    // Dev fallback: build from source if Go is available and source exists.
    const hasSource = fs.existsSync(path.join(PKG_ROOT, 'go.mod')) && fs.existsSync(path.join(PKG_ROOT, 'main.go'));
    if (hasSource) {
      const r = spawnSync('go', ['build', '-o', bin, '.'], { cwd: PKG_ROOT, stdio: 'inherit' });
      if (r.status === 0 && fs.existsSync(bin)) return bin;
    }
    const msg = `croctui binary not found.\nTried: ${bin}\nDownload: ${url}\n${e.message}\n` +
      `Fix: wait for GitHub Release ${tag} with asset ${asset}, or run with Go: go run .`;
    throw new Error(msg);
  }
}

async function main() {
  let bin;
  try {
    bin = await ensureBin();
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
  const res = spawnSync(bin, process.argv.slice(2), { stdio: 'inherit' });
  process.exit(res.status === null ? 1 : res.status);
}

main();
