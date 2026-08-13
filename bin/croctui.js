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

function localBin() { return path.join(__dirname, binName()); }
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

// Robust fetch: buffer in memory, handle redirects, timeout, then atomic write.
// Avoids Windows ENOENT on rename when the file handle is still open.
function fetch(url, dest) {
  return new Promise((resolve, reject) => {
    const req = https.get(url, { headers: { 'User-Agent': 'croctui-installer' } }, res => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        fetch(res.headers.location, dest).then(resolve, reject);
        return;
      }
      if (res.statusCode !== 200) {
        res.resume();
        reject(new Error(`download failed: ${res.statusCode} ${res.statusMessage} for ${url}`));
        return;
      }
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

async function ensureBin() {
  const bin = findBin();
  if (fs.existsSync(bin)) return bin;

  let version = '0.1.0';
  try { version = require('../package.json').version; } catch {}
  const tag = `v${version}`;
  const asset = binName();
  const base = process.env.CROCTUI_RELEASE_BASE || 'https://github.com/asterxsk/croctui/releases/download';
  const url = `${base}/${tag}/${asset}`;
  const tmp = bin + '.tmp';

  try {
    await fetch(url, tmp);
    // Windows rename fails if target exists or is locked — remove first, fallback to copy
    try { await fs.promises.unlink(bin); } catch {}
    try {
      await fs.promises.rename(tmp, bin);
    } catch {
      await fs.promises.copyFile(tmp, bin);
      await fs.promises.unlink(tmp);
    }
    if (process.platform !== 'win32') await fs.promises.chmod(bin, 0o755);
    return bin;
  } catch (e) {
    try { await fs.promises.unlink(tmp); } catch {}
    const hasSource = fs.existsSync(path.join(PKG_ROOT, 'go.mod')) && fs.existsSync(path.join(PKG_ROOT, 'main.go'));
    if (hasSource) {
      const r = spawnSync('go', ['build', '-o', bin, '.'], { cwd: PKG_ROOT, stdio: 'inherit' });
      if (r.status === 0 && fs.existsSync(bin)) return bin;
    }
    throw new Error(
      `croctui binary not found.\nTried: ${bin}\nDownload: ${url}\n${e.message}\n` +
      `Fix: npx will download ${asset} from ${tag} — check release exists, or run with Go: go run .`
    );
  }
}

async function main() {
  let bin;
  try { bin = await ensureBin(); } catch (e) { console.error(e.message); process.exit(1); }
  const res = spawnSync(bin, process.argv.slice(2), { stdio: 'inherit' });
  process.exit(res.status === null ? 1 : res.status);
}
main();
