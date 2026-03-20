<script lang="ts">
  let ex4File: File | null = $state(null);
  let csvFile: File | null = $state(null);
  let status: string = $state("");
  let isBacktesting: boolean = $state(false);
  let result: any = $state(null);

  function handleFile(event: Event, type: 'ex4' | 'csv') {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      if (type === 'ex4') ex4File = target.files[0];
      if (type === 'csv') csvFile = target.files[0];
    }
  }

  async function uploadAndRun() {
    if (!ex4File || !csvFile) {
      status = "Please upload both .ex4 and .csv files.";
      return;
    }
    
    isBacktesting = true;
    status = "Uploading and starting backtest...";
    result = null;

    const formData = new FormData();
    formData.append("ex4", ex4File);
    formData.append("history", csvFile);

    try {
      const response = await fetch("/api/backtest", {
        method: "POST",
        body: formData
      });
      
      const data = await response.json();
      
      if (response.ok) {
        status = "Backtest completed!";
        result = data;
      } else {
        status = "Error: " + (data.error || "Unknown error occurred.");
      }
    } catch (err) {
      status = "Request failed. Is the backend running?";
      console.error(err);
    } finally {
      isBacktesting = false;
    }
  }
</script>

<main class="container mx-auto p-8 max-w-2xl">
  <div class="bg-gray-800 rounded-xl shadow-2xl overflow-hidden border border-gray-700">
    <div class="p-6 border-b border-gray-700 bg-gray-900">
      <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
        MetaTrader 4 Headless Backtester
      </h1>
      <p class="text-sm text-gray-400 mt-2">Upload your Expert Advisor and History Data</p>
    </div>

    <div class="p-6 space-y-6">
      <!-- Upload EX4 -->
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-2">Expert Advisor (.ex4)</label>
        <div class="flex items-center justify-center w-full">
          <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-gray-600 border-dashed rounded-lg cursor-pointer bg-gray-700 hover:bg-gray-600 transition-colors">
            <div class="flex flex-col items-center justify-center pt-5 pb-6">
              <svg class="w-8 h-8 mb-4 text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"/>
              </svg>
              <p class="mb-2 text-sm text-gray-300"><span class="font-semibold">Click to upload</span> or drag and drop</p>
              <p class="text-xs text-gray-400">{ex4File ? ex4File.name : "No file selected"}</p>
            </div>
            <input type="file" class="hidden" accept=".ex4" onchange={(e) => handleFile(e, 'ex4')} />
          </label>
        </div>
      </div>

      <!-- Upload CSV -->
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-2">History Data (.csv)</label>
        <div class="flex items-center justify-center w-full">
          <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-gray-600 border-dashed rounded-lg cursor-pointer bg-gray-700 hover:bg-gray-600 transition-colors">
            <div class="flex flex-col items-center justify-center pt-5 pb-6">
              <svg class="w-8 h-8 mb-4 text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"/>
              </svg>
              <p class="mb-2 text-sm text-gray-300"><span class="font-semibold">Click to upload</span> or drag and drop</p>
              <p class="text-xs text-gray-400">{csvFile ? csvFile.name : "No file selected"}</p>
            </div>
            <input type="file" class="hidden" accept=".csv" onchange={(e) => handleFile(e, 'csv')} />
          </label>
        </div>
      </div>

      <!-- Action Button -->
      <button 
        onclick={uploadAndRun} 
        disabled={isBacktesting}
        class="w-full text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-800 disabled:cursor-not-allowed focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-3 text-center transition-colors"
      >
        {isBacktesting ? 'Running Backtest...' : 'Run Backtest'}
      </button>

      <!-- Status / Message -->
      {#if status}
        <div class="p-4 mb-4 text-sm rounded-lg border {status.includes('Error') || status.includes('failed') ? 'bg-red-900 border-red-500 text-red-200' : 'bg-blue-900 border-blue-500 text-blue-200'}" role="alert">
          {status}
        </div>
      {/if}

      <!-- Results Display -->
      {#if result}
        <div class="mt-6 p-4 bg-gray-900 rounded-lg border border-green-500/30">
          <h3 class="text-lg font-semibold text-green-400 mb-2">Results</h3>
          <pre class="text-xs text-gray-300 overflow-x-auto p-2 bg-black/50 rounded">{JSON.stringify(result, null, 2)}</pre>
        </div>
      {/if}
    </div>
  </div>
</main>
