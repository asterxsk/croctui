#!/usr/bin/env node
'use strict';

const { spawnSync } = require('node:child_process');
const fs = require('node:fs');
const path = require('node:path');

// Resolve paths relative to this package root, not the caller's cwd.
const PKG_ROOT = path.resolve(__dirname, '..');
const EXE = process.platform === 'win32' ? 'crocgui.exe' : 'crocgui';
const BIN = path.join(PKG_ROOT, 'bin', EXE);

// Build the Go binary if it's missing (or older than the source).
function needsBuild() {
  if (!fs.existsSync(BIN)) return true;
  const binMtime = fs.statSync(BIN).mtimeMs;
  const sources = ['main.go', 'go.mod', 'go.sum', 'ui', 'croc', 'dialog'];
  for (const s of sources) {
    const p = path.join(PKG_ROOT, s);
    if (fs.existsSync(p) && fs.statSync(p).mtimeMs > binMtime) return true;
  }
  return false;
}

function build() {
  const res = spawnSync('go', ['build', '-o', BIN, '.'], {
    cwd: PKG_ROOT,
    stdio: 'inherit',
  });
  if (res.status !== 0) {
    process.exit(res.status || 1);
  }
}

if (needsBuild()) {
  build();
}

// Run the binary, inheriting stdio and forwarding args.
const res = spawnSync(BIN, process.argv.slice(2), {
  cwd: PKG_ROOT,
  stdio: 'inherit',
});
process.exit(res.status === null ? 1 : res.status);
