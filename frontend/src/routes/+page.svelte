<script lang="ts">
	let url = $state('');
	let shortenedUrl = $state('');
	let errorMessage = $state('');
	let isLoading = $state(false);
	let isCopied = $state(false);

	let isInputValid = $derived(url.trim().length > 0 && url.startsWith('http'));

	async function shortenUrl(url: string): Promise<string | null> {
		const apiUrl = 'http://localhost:8080/shorten'; // 👈 REPLACE WITH YOUR BACKEND API URL

		try {
			const response = await fetch(apiUrl, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ url })
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to shorten URL on the server.');
			}

			const data = await response.json();

			return data['short_url'];
		} catch (error) {
			console.error('API call failed:', error);
			if (error instanceof Error) {
				throw new Error(`Error: ${error.message}`);
			}
			return null;
		}
	}

	// Handle the form submission
	async function handleSubmit() {
		errorMessage = '';
		shortenedUrl = '';
		isCopied = false;
		isLoading = true;

		try {
			if (!isInputValid) {
				throw new Error('Please enter a valid URL (e.g., https://example.com)');
			}

			const result = await shortenUrl(url);

			if (result) {
				shortenedUrl = result;
			} else {
				errorMessage = 'Failed to shorten URL. Please try again.';
			}
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
			} else {
				errorMessage = 'An unexpected error occurred.';
			}
		} finally {
			isLoading = false;
		}
	}

	// Handle copying the shortened URL to the clipboard
	function handleCopy() {
		const el = document.createElement('textarea');
		el.value = shortenedUrl;
		document.body.appendChild(el);
		el.select();
		document.execCommand('copy');
		document.body.removeChild(el);
		isCopied = true;

		setTimeout(() => {
			isCopied = false;
		}, 2000);
	}
</script>

<main class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
	<div class="w-full max-w-xl rounded-xl bg-white p-8 text-center shadow-2xl">
		<h1 class="mb-2 text-4xl font-bold text-gray-800">URL Shortener</h1>
		<p class="mb-8 text-gray-600">Shorten your long links into tiny, shareable ones.</p>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="flex flex-col gap-4"
		>
			<input
				type="url"
				bind:value={url}
				placeholder="Enter your URL here"
				class="w-full rounded-lg border border-gray-300 px-5 py-3 transition duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>
			<button
				type="submit"
				disabled={!isInputValid || isLoading}
				class="rounded-lg bg-blue-600 px-6 py-3 font-semibold text-white shadow-md transition duration-200 hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
			>
				{#if isLoading}
					Shortening...
				{:else}
					Shorten URL
				{/if}
			</button>
		</form>

		<div class="mt-6 min-h-24">
			{#if shortenedUrl}
				<div
					class="flex flex-col items-center justify-between rounded-lg bg-green-100 p-4 shadow-inner sm:flex-row"
				>
					<a
						href={'https://' + shortenedUrl}
						target="_blank"
						class="mb-2 break-all text-left font-medium text-green-800 sm:mb-0"
					>
						{shortenedUrl}
					</a>
					<button
						onclick={handleCopy}
						class="rounded-lg bg-green-600 px-4 py-2 font-medium text-white transition duration-200 hover:bg-green-700"
					>
						{#if isCopied}
							Copied!
						{:else}
							Copy
						{/if}
					</button>
				</div>
			{/if}

			<!-- Error Message Display -->
			{#if errorMessage}
				<div class="mt-4 rounded-lg bg-red-100 p-4 font-medium text-red-800 shadow-inner">
					{errorMessage}
				</div>
			{/if}
		</div>
	</div>
</main>
