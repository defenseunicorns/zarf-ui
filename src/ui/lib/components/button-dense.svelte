<!-- 
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
 -->
<script lang="ts">
	import { Button, type ButtonProps } from '@ui';
	import { current_component } from 'svelte/internal';

	export let dense = true;

	interface $$Props extends ButtonProps {
		dense?: boolean;
	}

	$: computedClass = `${(dense && 'dense') || ''} ${$$restProps.class || ''}`;
</script>

<Button eventComponent={current_component} {...$$restProps} class={computedClass}>
	<slot />
</Button>

<style global>
	.box.button.dense {
		height: 30px !important;
		padding: 4px 10px !important;
		font-size: small !important;
	}

	/* 
	 * Fix bug in UUI (mdc-button) with transparent disabled colors in button 
	 * button-label background color for label should always be transparent
	 * https://github.com/defenseunicorns/UnicornUI/issues/204
	*/
	.button.mdc-button.dense:disabled > .button-label {
		background-color: transparent !important;
	}
</style>
