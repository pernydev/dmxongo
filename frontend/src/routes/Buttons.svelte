<script>
	import { onMount } from "svelte";

	let buttons = {};

	function click(button) {
        fetch(`http://localhost:8080/function/${button}`, {
            method: 'POST'
        });

        if (buttons[button].type === 'toggle') {
            buttons[button].running = !buttons[button].running;
        }
	}

    onMount(async () => {
        const response = await fetch('http://localhost:8080/function');
        const data = await response.json();
        console.log(data);
        buttons = data;
    });
</script>

<grid>
	{#each Object.entries(buttons) as [button, { running, type }]}
		<button {running} {type} on:click={() => click(button)}>
			{button}
		</button>
	{/each}
</grid>

<style>
	grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		grid-template-rows: repeat(3, 1fr);
		gap: 10px;
	}

	button {
		padding: 10px;
		padding-block: 2rem;
		font-size: 1.5rem;
		border: 2px solid blue;
	}

	button[type='basic'] {
		border: 2px solid red;
	}

	button[running='true'] {
		background: red;
	}
</style>
