// Minimart Application JavaScript

// Utility functions
const Minimart = {
    // Bitcoin amount formatting
    formatBitcoin: function(satoshis) {
        const amount = parseInt(satoshis);
        
        if (amount >= 10000000) { // >= 0.1 BTC
            const btc = amount / 100000000;
            return `${btc.toFixed(8)} BTC`;
        } else if (amount >= 100000) { // >= 1 mBTC
            const mbtc = amount / 100000;
            return `${mbtc.toFixed(3)} mBTC`;
        } else {
            return `${amount} sats`;
        }
    },

    // Time formatting
    formatTimeRemaining: function(targetTime) {
        const now = new Date();
        const target = new Date(targetTime);
        const diff = target - now;
        
        if (diff <= 0) return "Ready now";
        
        const minutes = Math.ceil(diff / (1000 * 60));
        if (minutes < 60) {
            return `${minutes} min${minutes !== 1 ? 's' : ''}`;
        }
        
        const hours = Math.floor(minutes / 60);
        const remainingMins = minutes % 60;
        return `${hours}h ${remainingMins}m`;
    },

    // Toast notification helper
    showToast: function(message, type = 'info') {
        if (window.Datastar) {
            Datastar.store.signals.toast = {
                show: true,
                message: message,
                type: type
            };
            
            // Auto-hide after 5 seconds
            setTimeout(() => {
                if (Datastar.store.signals.toast.show) {
                    Datastar.store.signals.toast.show = false;
                }
            }, 5000);
        } else {
            // Fallback to basic alert if Datastar isn't available
            alert(message);
        }
    },

    // Loading state management
    setLoading: function(element, loading = true) {
        if (loading) {
            element.classList.add('loading');
            element.disabled = true;
        } else {
            element.classList.remove('loading');
            element.disabled = false;
        }
    },

    // Form validation helpers
    validateRequired: function(value) {
        return value && value.trim().length > 0;
    },

    validateEmail: function(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    },

    validateBitcoinAmount: function(amount) {
        const num = parseFloat(amount);
        return !isNaN(num) && num > 0;
    },

    // Bitcoin input helpers
    convertToBTC: function(satoshis) {
        return (parseInt(satoshis) / 100000000).toFixed(8);
    },

    convertToSatoshis: function(btc) {
        return Math.round(parseFloat(btc) * 100000000);
    },

    // Status badge helper
    getStatusBadgeClass: function(status) {
        const statusClasses = {
            'PENDING': 'status-pending',
            'ACCEPTED': 'status-accepted',
            'PREPARING': 'status-preparing',
            'READY': 'status-ready',
            'OUT_FOR_DELIVERY': 'status-out-for-delivery',
            'COMPLETED': 'status-completed',
            'CANCELLED': 'status-cancelled',
            'REJECTED': 'status-rejected'
        };
        return statusClasses[status] || 'status-pending';
    },

    // Audio notification (for merchants)
    playNotification: function() {
        // Create a simple beep sound
        const audioContext = new (window.AudioContext || window.webkitAudioContext)();
        const oscillator = audioContext.createOscillator();
        const gainNode = audioContext.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(audioContext.destination);
        
        oscillator.frequency.value = 800;
        oscillator.type = 'sine';
        
        gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
        gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5);
        
        oscillator.start(audioContext.currentTime);
        oscillator.stop(audioContext.currentTime + 0.5);
    },

    // Copy to clipboard helper
    copyToClipboard: function(text) {
        navigator.clipboard.writeText(text).then(() => {
            this.showToast('Copied to clipboard!', 'success');
        }).catch(() => {
            this.showToast('Failed to copy', 'error');
        });
    }
};

// Datastar event handlers
document.addEventListener('DOMContentLoaded', function() {
    // Initialize any global state if Datastar is available
    if (window.Datastar) {
        // Add global helper functions to Datastar store
        Datastar.store.helpers = {
            formatBitcoin: Minimart.formatBitcoin,
            formatTimeRemaining: Minimart.formatTimeRemaining,
            getStatusBadgeClass: Minimart.getStatusBadgeClass
        };
    }

    // Handle form submissions
    document.addEventListener('submit', function(e) {
        const form = e.target;
        
        // Add loading state to submit buttons
        const submitBtn = form.querySelector('button[type="submit"]');
        if (submitBtn) {
            Minimart.setLoading(submitBtn, true);
            
            // Reset loading state after a delay (will be overridden by response)
            setTimeout(() => {
                Minimart.setLoading(submitBtn, false);
            }, 10000);
        }
    });

    // Handle Bitcoin amount inputs
    document.querySelectorAll('input[data-bitcoin-input]').forEach(function(input) {
        input.addEventListener('input', function() {
            const value = this.value;
            const preview = this.parentNode.querySelector('[data-bitcoin-preview]');
            
            if (preview && value) {
                const satoshis = parseFloat(value) * 100000000;
                preview.textContent = `≈ ${Minimart.formatBitcoin(satoshis)}`;
            }
        });
    });

    // Auto-refresh timestamps
    setInterval(function() {
        document.querySelectorAll('[data-timestamp]').forEach(function(element) {
            const timestamp = element.dataset.timestamp;
            const formatted = Minimart.formatTimeRemaining(timestamp);
            if (element.textContent !== formatted) {
                element.textContent = formatted;
            }
        });
    }, 30000); // Update every 30 seconds
});

// Service Worker registration for PWA (optional)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function() {
        navigator.serviceWorker.register('/sw.js').then(function(registration) {
            console.log('SW registered: ', registration);
        }).catch(function(registrationError) {
            console.log('SW registration failed: ', registrationError);
        });
    });
}

// Export for module usage
if (typeof module !== 'undefined' && module.exports) {
    module.exports = Minimart;
}

// Make available globally
window.Minimart = Minimart;
