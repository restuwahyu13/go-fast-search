import { useState, useEffect, useRef, useCallback } from 'react'

interface FormattedPost {
	name?: string
	email?: string
	phone?: string
	date_of_birth?: string
	age?: string
}

interface User {
	id: string
	name: string
	email: string
	phone: string
	date_of_birth: string
	age: number
	_formatted?: FormattedPost
}

interface ApiResponse {
	data?: {
		results: User[]
		total: number
	}
}

function formatHumanReadable(num: number) {
	if (num >= 1_000_000_000) {
		return `${(num / 1_000_000_000).toFixed(1).replace(/\.0$/, '')} miliar`
	}
	if (num >= 1_000_000) {
		return `${(num / 1_000_000).toFixed(1).replace(/\.0$/, '')} juta`
	}
	if (num >= 1_000) {
		return `${(num / 1_000).toFixed(1).replace(/\.0$/, '')} ribu`
	}
	return `${num} data`
}

export default function Home() {
	const [users, setUsers] = useState<User[]>([])
	const [searchTerm, setSearchTerm] = useState('')
	const [page, setPage] = useState(1)
	const [loading, setLoading] = useState(false)
	const [paginationLoading, setPaginationLoading] = useState(false)
	const [responseTime, setResponseTime] = useState<number | null>(null)
	const [error, setError] = useState<string | null>(null)
	const [hasMore, setHasMore] = useState(true)
	const observerRef = useRef<IntersectionObserver | null>(null)
	const loadMoreRef = useRef<HTMLTableRowElement | null>(null)
	const tableRef = useRef<HTMLDivElement | null>(null)
	const [total, setTotal] = useState<number>(0)

	const LIMIT = 10
	const MAX_USERS = 1000

	const fetchUsers = async (pageNum: number, isNewSearch: boolean = false) => {
		if (pageNum * LIMIT > MAX_USERS) {
			setHasMore(false)
			return
		}

		setLoading(isNewSearch)
		setPaginationLoading(!isNewSearch)
		setError(null)
		const startTime = performance.now()

		try {
			const url = searchTerm
				? `${process.env.NEXT_PUBLIC_API_URL}/users?limit=${LIMIT}&page=${pageNum}&search=${encodeURIComponent(searchTerm)}`
				: `${process.env.NEXT_PUBLIC_API_URL}/users?limit=${LIMIT}&page=${pageNum}`
			const response = await fetch(url)
			if (!response.ok) {
				throw new Error(`HTTP error! Status: ${response.status}`)
			}
			const data: ApiResponse = await response.json()
			const newUsers = data?.data?.results || []

			setUsers((prev) => (isNewSearch ? newUsers : [...prev, ...newUsers]))
			setHasMore(newUsers.length === LIMIT && users.length + newUsers.length < MAX_USERS)

			setTotal(data?.data?.total || 0)
			setResponseTime(performance.now() - startTime)
		} catch (error) {
			console.error('Error fetching users:', error)
			setError('Failed to fetch users. Please try again.')
			if (isNewSearch) setUsers([])
		} finally {
			setLoading(false)
			setPaginationLoading(false)
		}
	}

	useEffect(() => {
		setPage(1)
		setUsers([])
		setHasMore(true)
		const timeoutId = setTimeout(() => fetchUsers(1, true), 300)
		return () => clearTimeout(timeoutId)
	}, [searchTerm])

	const handleObserver = useCallback(
		(entries: IntersectionObserverEntry[]) => {
			const target = entries[0]
			if (target.isIntersecting && hasMore && !loading && !paginationLoading) {
				setPage((prev) => prev + 1)
			}
		},
		[hasMore, loading, paginationLoading]
	)

	useEffect(() => {
		if (page > 1) {
			fetchUsers(page)
		}
	}, [page])

	useEffect(() => {
		if (tableRef.current && loadMoreRef.current) {
			observerRef.current = new IntersectionObserver(handleObserver, {
				root: tableRef.current,
				rootMargin: '50px',
				threshold: 0.1
			})
			observerRef.current.observe(loadMoreRef.current)
		}

		return () => {
			if (observerRef.current && loadMoreRef.current) {
				observerRef.current.unobserve(loadMoreRef.current)
			}
		}
	}, [handleObserver])

	const renderHighlightedText = (text: string | undefined): React.ReactElement | null => {
		if (!text) return null
		return <span dangerouslySetInnerHTML={{ __html: text }} />
	}

	const formatDate = (date: string | undefined): string => {
		return date ? date.split('T')[0] : ''
	}

	return (
		<div className='min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4 sm:p-6'>
			<div className='max-w-7xl mx-auto'>
				{/* Header */}
				<h1 className='text-2xl sm:text-3xl font-extrabold text-center mb-6 text-indigo-700 tracking-tight'>Faster Search With Golang + MeiliSearch</h1>

				{/* Search Input */}
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
				</div>

				<div className='mt-2 space-y-1 flex flex-cols justify-between mb-6 p-2'>
					{responseTime !== null && <p className='text-sm text-gray-500 font-bold'>API Response Time: {responseTime.toFixed(2)} ms</p>}
					{responseTime !== null && <p className='text-sm text-gray-500 font-bold'>Total Data API: {formatHumanReadable(total)}</p>}
					{error && <p className='text-sm text-red-500'>{error}</p>}
				</div>

				{/* Users Table with Scroll */}
				<div className='bg-white shadow-lg rounded-xl overflow-hidden'>
					<div ref={tableRef} className='max-h-[590px] overflow-y-auto'>
						<table className='min-w-full divide-y divide-gray-200'>
							<thead className='bg-indigo-50 sticky top-0 z-10'>
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
										<td colSpan={5} className='px-4 py-6 sm:px-6 text-center text-sm text-gray-500'>
											<div className='flex justify-center items-center'>
												<svg className='animate-spin h-5 w-5 text-indigo-500' xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24'>
													<circle className='opacity-25' cx='12' cy='12' r='10' stroke='currentColor' strokeWidth='4' />
													<path className='opacity-75' fill='currentColor' d='M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z' />
												</svg>
												<span className='ml-2'>Loading...</span>
											</div>
										</td>
									</tr>
								) : users.length > 0 ? (
									<>
										{users.map((user) => (
											<tr key={user.id} className='hover:bg-indigo-50 transition-colors duration-200'>
												<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>
													{user._formatted?.name ? renderHighlightedText(user._formatted.name) : user.name}
												</td>
												<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>
													{user._formatted?.email ? renderHighlightedText(user._formatted.email) : user.email}
												</td>
												<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>
													{user._formatted?.phone ? renderHighlightedText(user._formatted.phone) : user.phone}
												</td>
												<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>
													{user._formatted?.date_of_birth ? renderHighlightedText(formatDate(user._formatted.date_of_birth)) : formatDate(user.date_of_birth)}
												</td>
												<td className='px-4 py-4 text-sm text-gray-900 sm:px-6 font-medium whitespace-normal break-words'>
													{user._formatted?.age ? renderHighlightedText(user._formatted.age) : user.age}
												</td>
											</tr>
										))}
										{hasMore && (
											<tr ref={loadMoreRef}>
												<td colSpan={5} className='px-4 py-4 sm:px-6'></td>
											</tr>
										)}
									</>
								) : (
									<tr>
										<td colSpan={5} className='px-4 py-6 sm:px-6 text-center text-sm text-gray-500'>
											No users found
										</td>
									</tr>
								)}
								{paginationLoading && (
									<tr>
										<td colSpan={5} className='px-4 py-4 sm:px-6 text-center text-sm text-gray-500'>
											<div className='flex justify-center items-center'>
												<svg className='animate-spin h-5 w-5 text-indigo-500' xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24'>
													<circle className='opacity-25' cx='12' cy='12' r='10' stroke='currentColor' strokeWidth='4' />
													<path className='opacity-75' fill='currentColor' d='M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z' />
												</svg>
												<span className='ml-2'>Loading more...</span>
											</div>
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
