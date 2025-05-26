import { useState, useEffect } from 'react'

interface Post {
	id: string
	name: string
	email: string
	phone: string
	date_of_birth: string
	age: number
	city: string
	state: string
	direction: string
	country: string
}

export default function Home() {
	const [users, setUsers] = useState<Post[]>([])
	const [searchTerm, setSearchTerm] = useState('')
	const [loading, setLoading] = useState(false)
	const [responseTime, setResponseTime] = useState<number | null>(null)

	useEffect(() => {
		const fetchUsers = async () => {
			setLoading(true)
			const startTime = performance.now()

			try {
				const url = searchTerm ? `http://localhost:4000/api/v1/users?limit=10&page=1&search=${encodeURIComponent(searchTerm)}` : 'http://localhost:4000/api/v1/users?limit=10&page=1'
				const response = await fetch(url)

				const data = await response.json()
				setUsers(data?.data?.results || [])

				const endTime = performance.now()
				setResponseTime(endTime - startTime)
			} catch (error) {
				console.error('Error fetching users:', error)
				setUsers([])
				setResponseTime(null)
			} finally {
				setLoading(false)
			}
		}

		const timeoutId = setTimeout(() => {
			fetchUsers()
		}, 300)

		return () => clearTimeout(timeoutId)
	}, [searchTerm])

	return (
		<div className='min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4 sm:p-6'>
			<div className='max-w-7xl mx-auto'>
				<h1 className='text-2xl sm:text-3xl font-extrabold text-center mb-6 text-indigo-700 tracking-tight'>Faster Search With Golang + MeiliSearch</h1>
				<div className='mb-6 relative'>
					<input
						type='text'
						placeholder='Search users...'
						value={searchTerm}
						onChange={(e) => setSearchTerm(e.target.value)}
						className='w-full p-4 pr-12 bg-white border border-gray-200 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all duration-300 text-gray-700 placeholder-gray-400'
					/>
					<svg className='absolute right-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400' fill='none' stroke='currentColor' viewBox='0 0 24 24' xmlns='http://www.w3.org/2000/svg'>
						<path strokeLinecap='round' strokeLinejoin='round' strokeWidth='2' d='M21 21l-4.35-4.35m1.85-5.15a7 7 0 11-14 0 7 7 0 0114 0z' />
					</svg>
					{responseTime !== null && <p className='mt-2 text-sm text-gray-500'>API Response Time: {responseTime.toFixed(2)} ms</p>}
				</div>

				<div className='bg-white shadow-lg rounded-xl overflow-hidden'>
					<div className='overflow-x-auto'>
						<table className='min-w-full divide-y divide-gray-200'>
							<thead className='bg-indigo-50'>
								<tr>
									<th className='px-4 py-3 text-left text-xs font-semibold text-indigo-700 uppercase tracking-wider sm:px-6 min-w-[120px] sm:min-w-[150px]'>Name</th>
									<th className='px-4 py-3 text-left text-xs font-semibold text-indigo-700 uppercase tracking-wider sm:px-6 min-w-[180px] sm:min-w-[200px]'>Email</th>
									<th className='px-4 py-3 text-left text-xs font-semibold text-indigo-700 uppercase tracking-wider sm:px-6 min-w-[140px] sm:min-w-[160px]'>Phone</th>
									<th className='px-4 py-3 text-left text-xs font-semibold text-indigo-700 uppercase tracking-wider sm:px-6 min-w-[120px] sm:min-w-[140px]'>Date of Birth</th>
									<th className='px-4 py-3 text-left text-xs font-semibold text-indigo-700 uppercase tracking-wider sm:px-6 min-w-[80px] sm:min-w-[100px]'>Age</th>
								</tr>
							</thead>

							<tbody className='bg-white divide-y divide-gray-200'>
								{loading ? (
									<tr>
										<td colSpan={9} className='px-4 py-6 sm:px-6 text-center text-sm text-gray-500'>
											<div className='flex justify-center items-center'>
												<svg className='animate-spin h-5 w-5 text-indigo-500' xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24'>
													<circle className='opacity-25' cx='12' cy='12' r='10' stroke='currentColor' strokeWidth='4' />
													<path className='opacity-75' fill='currentColor' d='M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z' />
												</svg>
												<span className='ml-2'>Loading...</span>
											</div>
										</td>
									</tr>
								) : users?.length > 0 ? (
									users?.map((post) => (
										<tr key={post.id} className='hover:bg-indigo-50 transition-colors duration-200'>
											<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>{post?.name}</td>
											<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>{post?.email}</td>
											<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>{post?.phone}</td>
											<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>{post?.date_of_birth?.split('T')[0]}</td>
											<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>{post?.age}</td>
										</tr>
									))
								) : (
									<tr>
										<td colSpan={9} className='px-4 py-6 sm:px-6 text-center text-sm text-gray-500'>
											No users found
										</td>
									</tr>
								)}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</div>
	)
}
