<!-- 
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
 -->
<script lang="ts">
	import { page } from '$app/stores';
	import { Packages } from '$lib/api';
	import { pkgSbomStore } from '$lib/store';
	import { Box, Stepper, type SSX } from '@ui';
	import { onDestroy } from 'svelte';

	const stepMap: Record<string, number> = {
		packages: 1,
		configure: 2,
		review: 3,
		deploy: 4,
	};

	// Clean any SBOMs that may have been extracted.
	onDestroy(() => {
		pkgSbomStore.set(undefined);
		Packages.cleanSBOM();
	});

	$: stepName = $page.route.id?.split('/').pop() || '';
	$: stepNumber = (stepName && stepMap[stepName]) || 500;
	$: getIconContent = (number: number): string | undefined =>
		stepNumber <= number ? number.toString() : undefined;
	$: stepDisabled = (number: number): boolean => (stepNumber < number ? true : false);

	const ssx: SSX = {
		$self: {
			display: 'flex',
			flexDirection: 'column',
			gap: '32px',
		},

		$xs: {
			$self: {
				padding: '0px 0px',
			},
		},
		$lg: {
			$self: {
				padding: '0px 2rem',
			},
		},
		$xl: {
			$self: {
				padding: '0px 16rem',
			},
		},
	};
</script>

<Box {ssx} class="packages-page">
	<div class="deploy-stepper-container">
		<Stepper
			color={'text-secondary-on-background'}
			orientation="horizontal"
			steps={[
				{
					title: 'Select',
					iconContent: getIconContent(stepMap.packages),
					variant: 'primary',
				},

				{
					title: 'Configure',
					iconContent: getIconContent(stepMap.configure),
					disabled: stepDisabled(stepMap.configure),
					variant: 'primary',
				},
				{
					title: 'Review',
					iconContent: getIconContent(stepMap.review),
					disabled: stepDisabled(stepMap.review),
					variant: 'primary',
				},
				{
					title: 'Deploy',
					iconContent: '4',
					disabled: stepDisabled(stepMap.deploy),
					variant: 'primary',
				},
			]}
		/>
	</div>
	<slot />
</Box>

<style>
	.deploy-stepper-container {
		max-width: 600px;
		margin: 0 auto;
		width: 100%;
	}
</style>
