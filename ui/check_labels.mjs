import { chromium } from 'playwright';

const browser = await chromium.launch({ headless: true, args: ['--ignore-certificate-errors'] });
const page = await browser.newPage();

// Go to login page
await page.goto('http://localhost:8080/login', { timeout: 15000 });
await page.waitForTimeout(1000);

await page.fill('input[id="username"]', 'admin');
await page.fill('input[type="password"]', 'admin');
await page.click('button[type="submit"]');
await page.waitForTimeout(2000);

console.log('=== URL after login:', page.url());

// Go to scopes
await page.goto('http://localhost:8080/scopes', { timeout: 15000 });
await page.waitForTimeout(3000);
await page.screenshot({ path: '/tmp/scopes-kind.png', fullPage: true });
console.log('=== SCOPES URL:', page.url());

const body = await page.textContent('body');
console.log('=== PAGE BODY (first 3000 chars):');
console.log(body.substring(0, 3000));

// Also check API
const apiRes = await (await page.request).get('http://localhost:8080/api/v1/scopes');
const apiData = await apiRes.json();
console.log(`=== API returned ${apiData.length} scopes`);
for (const s of apiData) {
  console.log(`  ${s.metadata.name}: labels=`, JSON.stringify(s.metadata.labels));
}

await browser.close();
