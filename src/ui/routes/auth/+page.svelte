<!-- 
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
 -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Auth } from '$lib/api';
	import { Hero } from '$lib/components';
	import { Typography } from '@ui';
	import sadDay from '@images/sadness.png';

	let authFailure = false;

	page.subscribe(async ({ url }) => {
		let token = url.searchParams.get('token') || '';
		let next = url.searchParams.get('next');

		if (await Auth.connect(token)) {
			goto(next || '/');
		} else {
			authFailure = true;
		}
	});
</script>

{#if authFailure}
	<Hero>
		<img src={sadDay} alt="Sadness" id="sadness" width="40%" />

		<Typography variant="h5">Could not authenticate!</Typography>
		<Typography variant="body2">
			Please make sure you are using the complete link to connect provided by Zarf.
		</Typography>
	</Hero>
{/if}
