<!-- 
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
 -->
<script lang="ts">
	import { Paper, Typography, type PaperProps } from '@ui';
	import { current_component } from 'svelte/internal';

	type T = $$Generic<EventTarget>;
	export let backgroundColor = 'chip-background-color';
	export let color = 'chip-color';
	export let ripple = false;

	interface $$Props extends PaperProps<T> {
		ripple?: boolean;
	}

	$: computedClass = `zarf-chip ${(ripple && 'ripple') || ''} ${$$restProps.class || ''}`;
</script>

<Paper
	eventComponent={current_component}
	{backgroundColor}
	{color}
	{...$$restProps}
	class={computedClass}
>
	<slot name="leading-icon" />
	<Typography element="span" variant="zarf-chip-typography">
		<slot />
	</Typography>
	<slot name="trailing-icon" />
</Paper>

<style lang="scss" global>
	.zarf-chip.paper {
		border-radius: 16px;
		display: flex;
		align-items: center;
		justify-content: center;
		word-break: break-all;
		padding: 4px 8px;
		gap: 1px;
	}
</style>
