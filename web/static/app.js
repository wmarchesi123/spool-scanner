// Copyright 2025 William Marchesi

// Author: William Marchesi
// Email: will@marchesi.io
// Website: https://marchesi.io/

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

document.addEventListener('alpine:init', () => {
    Alpine.data('spoolSelector', () => ({
        // State
        loading: true,
        error: null,
        step: 'loading', // 'loading', 'select-printer', 'confirming', 'success', 'error'
        
        // Data
        spoolId: SPOOL_ID,
        spool: null,
        printers: [],
        selectedPrinter: null,
        
        // Methods
        async init() {
            console.log('Initializing spool selector for spool:', this.spoolId);
            
            try {
                // Load spool data and printers in parallel
                const [spoolData, printersData] = await Promise.all([
                    this.fetchSpool(),
                    this.fetchPrinters()
                ]);
                
                this.spool = spoolData;
                this.printers = printersData.printers;
                
                // If only one printer is ready, pre-select it
                const readyPrinters = this.printers.filter(p => p.status === 'Ready');
                if (readyPrinters.length === 1) {
                    this.selectedPrinter = readyPrinters[0];
                }
                
                this.loading = false;
                this.step = 'select-printer';
            } catch (err) {
                console.error('Initialization error:', err);
                this.error = err.message || 'Failed to load data';
                this.step = 'error';
                this.loading = false;
            }
        },
        
        async fetchSpool() {
            const response = await fetch(`/api/spool/${this.spoolId}`);
            if (!response.ok) {
                throw new Error('Spool not found');
            }
            return await response.json();
        },
        
        async fetchPrinters() {
            const response = await fetch('/api/printers');
            if (!response.ok) {
                throw new Error('Failed to load printers');
            }
            return await response.json();
        },
        
        selectPrinter(printer) {
            this.selectedPrinter = printer;
        },
        
        async confirmAssignment() {
            if (!this.selectedPrinter) {
                this.error = 'Please select a printer';
                return;
            }
            
            this.step = 'confirming';
            this.error = null;
            
            try {
                const response = await fetch('/api/assign', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        spool_id: this.spoolId,
                        printer_id: this.selectedPrinter.id
                    })
                });
                
                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || 'Failed to assign spool');
                }
                
                this.step = 'success';
                
                // Redirect to the printer's OctoPrint URL after 3 seconds
                setTimeout(() => {
                    if (this.selectedPrinter.url) {
                        window.location.href = this.selectedPrinter.url;
                    } else {
                        window.location.href = '/';
                    }
                }, 3000);
                
            } catch (err) {
                console.error('Assignment error:', err);
                this.error = err.message || 'Failed to assign spool';
                this.step = 'select-printer';
            }
        },
        
        formatWeight(grams) {
            return Math.round(grams) + ' g';
        },
        
        getPrinterStatusColor(status) {
            switch (status) {
                case 'Ready':
                    return 'green';
                case 'Printing':
                    return 'orange';
                case 'Error':
                    return 'red';
                default:
                    return 'gray';
            }
        },
        
        getPrinterStatusEmoji(status) {
            switch (status) {
                case 'Ready':
                    return '✅';
                case 'Printing':
                    return '🖨️';
                case 'Error':
                    return '❌';
                default:
                    return '❓';
            }
        }
    }));
});