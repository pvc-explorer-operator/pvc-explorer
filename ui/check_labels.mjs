import { chromium } from 'playwright';

const browser = await chromium.launch({ headless: true, args: ['--ignore-certificate-errors'] });
const page = await browser.newPage();

// Go to login page
await page.goto('https://pvc-explorer.hafas-analytics.internal/login', { timeout: 15000 });
await page.waitForTimeout(1000);

// Dump login page HTML
const loginHtml = await page.content();
console.log('=== LOGIN HTML snippet:');
console.log(loginHtml.substring(0, 3000));

// Fill login form
const inputs = await page.$$('input');
console.log(`Found ${inputs.length} inputs`);
for (const input of inputs) {
  const placeholder = await input.getAttribute('placeholder');
  const type = await input.getAttribute('type');
  const id = await input.getAttribute('id');
  console.log(`  input: placeholder="${placeholder}" type="${type}" id="${id}"`);
}

// Try various selectors
const usernameInput = await page.$('input[type="text"], input[name="username"], input[placeholder="Username"], input[id="username"]');
if (usernameInput) {
  await usernameInput.fill('admin');
  console.log('Filled username');
}

const passwordInput = await page.$('input[type="password"], input[name="password"], input[placeholder="Password"]');
if (passwordInput) {
  await passwordInput.fill('admin');
  console.log('Filled password');
}

const submitBtn = await page.$('button[type="submit"], button:has-text("Sign In"), button:has-text("Login")');
if (submitBtn) {
  await submitBtn.click();
  console.log('Clicked submit');
  await page.waitForTimeout(3000);
}

console.log('=== URL after login:', page.url());

// Go to scopes
await page.goto('https://pvc-explorer.hafas-analytics.internal/scopes', { timeout: 15000 });
await page.waitForTimeout(3000);
await page.screenshot({ path: '/tmp/scopes.png', fullPage: true });
console.log('=== SCOPES URL:', page.url());

const body = await page.textContent('body');
console.log('=== PAGE BODY (first 2000 chars):');
console.log(body.substring(0, 2000));

await browser.close();
