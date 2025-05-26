module.exports = {
	apps: [
		{
			name: 'web',
			script: 'npm',
			args: 'run start',
			watch: false,
			interpreter: 'none',
			max_memory_restart: '512M',
			listen_timeout: 3000,
			kill_timeout: 6000,
			combine_logs: true
		}
	]
}