<script>
	import { onMount } from 'svelte';
	import SpeedGraph from './SpeedGraph.svelte';
	import Buttons from './Buttons.svelte';

	let fixtures = [];
	let runFunctionDialog;

	onMount(() => {
		const ws = new WebSocket('ws://localhost:8080/ws/fixtures');
		ws.onmessage = (event) => {
			fixtures = JSON.parse(event.data);
		};
	});
</script>

<fixtures>
	{#each fixtures as fixture}
		<fixture
			style="--color: rgb({fixture.color.red}, {fixture.color.green}, {fixture.color.blue}); --brightness: {fixture.brightness / 255}"
		/>
	{/each}
</fixtures>

<SpeedGraph />

<dialog bind:this={runFunctionDialog}>
	<h1>Run Function</h1>
	<label>
		Function Name
		<input type="text" name="functionName" />
	</label>
</dialog>

<Buttons />

<style>
	fixtures {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}

	fixture {
		--color: rgb(0, 0, 0);
		background-color: var(--color);
		width: 100px;
		height: 100px;
		display: block;
		border-radius: 9999px;

		filter: brightness(var(--brightness));
	}
</style>
