import fs from "node:fs/promises";
import path from "node:path";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { chromium } = require("playwright");

const baseUrl = process.env.SMARTTODO_DEMO_URL || "http://127.0.0.1:4173/";
const outputDir = process.env.SMARTTODO_SMOKE_OUTPUT || "artifacts/smarttodo-demo-smoke";
const apiKey = process.env.SMARTTODO_SMOKE_API_KEY || "sk-demo-smoke";

await fs.mkdir(outputDir, { recursive: true });

const browser = await chromium.launch({ headless: true });
const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
const logs = [];

page.on("console", (msg) => logs.push(`[console:${msg.type()}] ${msg.text()}`));
page.on("pageerror", (err) => logs.push(`[pageerror] ${err.message}`));

try {
  await page.goto(baseUrl, { waitUntil: "networkidle" });
  await page.waitForSelector("#wasm-status");
  await page.waitForFunction(() => document.querySelector("#wasm-status")?.textContent?.includes("Ready"));

  await page.fill("#api-key-input", apiKey);
  await page.click("#connect-key");
  await page.waitForFunction(() => document.querySelector("#key-gate")?.classList.contains("hidden"));

  await page.locator("#terminal-root").click();
  await page.keyboard.press("Enter");
  await page.waitForTimeout(600);
  await page.keyboard.type("CI User");
  await page.keyboard.press("Enter");
  await page.waitForTimeout(600);
  await page.keyboard.type("CI Board");
  await page.keyboard.press("Enter");
  await page.waitForTimeout(1500);

  await page.screenshot({
    path: path.join(outputDir, "smarttodo-demo-smoke.png"),
    fullPage: true,
  });

  const sawBoot = logs.some((line) => line.includes("[smarttodo-go] boot"));
  const sawConnect = logs.some((line) => line.includes("[smarttodo-go] connect TUI connected"));
  if (!sawBoot || !sawConnect) {
    throw new Error("Go/WASM TUI boot logs were not observed in the browser console");
  }

  const statusText = await page.locator("#wasm-status").textContent();
  if (!statusText || !statusText.includes("Ready")) {
    throw new Error(`Unexpected runtime status: ${statusText ?? "<empty>"}`);
  }
} finally {
  await fs.writeFile(path.join(outputDir, "console.log"), logs.join("\n"), "utf8");
  await browser.close();
}
