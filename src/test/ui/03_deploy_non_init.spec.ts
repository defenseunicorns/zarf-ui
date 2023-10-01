// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

import { expect, test } from '@playwright/test';

test.beforeEach(async ({ page }) => {
	page.on('pageerror', (err) => console.log(err.message));
});

const getToSelectPage = async (page) => {
	await page.goto('/auth?token=insecure&next=/packages', { waitUntil: 'networkidle' });
};

const getToReview = async (page) => {
	await getToSelectPage(page);
	const expanded = (await page.locator('.button-label:has-text("Search Directory")')).first();
	if (expanded.isVisible()) {
		await expanded.click();
	}
	// Find first dos-games package deploy button.
	const dosGames = page.getByTitle('dos-games').first();
	// click the dos-games package deploy button.
	await dosGames.click();
	await page.getByRole('link', { name: 'review deployment' }).click();
	await page.waitForURL('/packages/dos-games/review');
};

test('deploy the dos-games package @post-init', async ({ page }) => {
	await getToReview(page);

	// Validate that the SBOM has been loaded
	const sbomInfo = await page.waitForSelector('#sbom-info', { timeout: 20000 });
	expect(await sbomInfo.innerText()).toMatch(/[0-9]+ artifacts to be reviewed/);

	await page.getByRole('link', { name: 'deploy package' }).click();
	await page.waitForURL('/packages/dos-games/deploy', { waitUntil: 'networkidle' });

	// verify the deployment succeeded
	await expect(page.locator('text=Deployment Succeeded')).toBeVisible({ timeout: 120000 });

	// then verify the page redirects to the Landing Page
	await page.waitForURL('/', { timeout: 10000 });
});
