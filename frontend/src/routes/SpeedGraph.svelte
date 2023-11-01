<script>
	import { onMount } from 'svelte';
	import { LinkedChart, LinkedLabel, LinkedValue } from 'svelte-tiny-linked-charts';

    const data = {};

    let max = 0;
    let min = 0;

    onMount(async () => {
        const ws = new WebSocket('ws://localhost:8080/ws/stats');
        ws.onmessage = (event) => {
            const json = JSON.parse(event.data).updateSpeeds;
            for (const key in json) {
                data[key] = (json[key] / 1000000).toFixed(2);
            }
        };
    });

    $: max = Math.max(...Object.values(data));
    $: min = Math.min(...Object.values(data));
</script>

<LinkedChart { data } showValue valueDefault="0" type="line" />


<small>
    MAX {max} ms
</small>

<small>
    MIN {min} ms
</small>
