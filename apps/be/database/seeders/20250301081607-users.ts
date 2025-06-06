require('dotenv/config')
import os from 'node:os'
import { randomUUID } from 'node:crypto'
import { QueryInterface, Sequelize } from 'sequelize'
import { faker } from '@faker-js/faker'
import { Index, MeiliSearch, RecordAny } from 'meilisearch'

process.env.NODE_OPTIONS = '--max-old-space-size=8192'
process.env.UV_THREADPOOL_SIZE = `${os.cpus().length}`

module.exports = {
	up: async (queryInterface: QueryInterface, _sequelize: Sequelize) => {
		const users: Record<string, any>[] = [
			{
				id: randomUUID(),
				name: 'onta cavalera',
				email: 'onta_cavalera@gmail.com',
				phone: faker.phone.number({ style: 'international' }),
				date_of_birth: faker.date.birthdate({ refDate: '1980-01-01' })?.toISOString(),
				age: `${faker.number.int({ min: 18, max: 200 })}`,
				address: faker.location.streetAddress(),
				city: faker.location.city(),
				state: faker.location.state(),
				direction: faker.location.direction(),
				country: faker.location.country(),
				postal_code: faker.location.zipCode(),
				created_at: new Date().toISOString(),
				updated_at: null,
				deleted_at: null
			},
			{
				id: randomUUID(),
				name: 'yajid laura',
				email: 'yajidlaura@gmail.com',
				phone: faker.phone.number({ style: 'international' }),
				date_of_birth: faker.date.birthdate({ refDate: '1980-01-01' })?.toISOString(),
				age: `${faker.number.int({ min: 18, max: 200 })}`,
				address: faker.location.streetAddress(),
				city: faker.location.city(),
				state: faker.location.state(),
				direction: faker.location.direction(),
				country: faker.location.country(),
				postal_code: faker.location.zipCode(),
				created_at: new Date().toISOString(),
				updated_at: null,
				deleted_at: null
			}
		]
		const totalData: number = 1000000

		for (let i = 1; i <= totalData; i++) {
			users.push({
				id: randomUUID(),
				name: faker.person.fullName(),
				email: faker.internet.email(),
				phone: faker.phone.number({ style: 'international' }),
				date_of_birth: faker.date.birthdate({ refDate: '1980-01-01' })?.toISOString(),
				age: `${faker.number.int({ min: 18, max: 200 })}`,
				address: faker.location.streetAddress(),
				city: faker.location.city(),
				state: faker.location.state(),
				direction: faker.location.direction(),
				country: faker.location.country(),
				postal_code: faker.location.zipCode(),
				created_at: new Date().toISOString(),
				updated_at: null,
				deleted_at: null
			})
		}

		console.log(`Total data users: ${users.length}`)

		if (users?.length >= totalData) {
			const client: MeiliSearch = new MeiliSearch({
				host: process.env.MEILI_DSN,
				apiKey: process.env.MEILI_MASTER_KEY
			})

			const newUsers: Record<string, any>[] = users.map((user: Record<string, any>) => ({
				id: user?.id,
				name: user?.name,
				email: user?.email,
				phone: user?.phone,
				date_of_birth: `${user?.date_of_birth}`,
				age: user?.age || '20',
				address: user?.address,
				city: user?.city,
				state: user?.state,
				direction: user?.direction,
				country: user?.country,
				postal_code: user?.postal_code,
				created_at: Math.floor(Date.now() / 1000),
				updated_at: null,
				deleted_at: null
			}))

			const index: Index<RecordAny> = client.index('users')
			index.updateSettings({
				pagination: { maxTotalHits: 1000 },
				searchableAttributes: ['name', 'email', 'phone'],
				sortableAttributes: ['created_at'],
				filterableAttributes: ['id', 'age', 'state', 'direction', 'country', 'created_at', 'deleted_at']
			})

			return await index
				.addDocuments(newUsers, { primaryKey: 'id' })
				.then(() => {
					console.log('MeiliSearch inserted successfully')
					return queryInterface
						.bulkInsert('users', users, { logging: true })
						.then(() => console.log('Sequelize inserted successfully'))
						.catch(console.error)
				})
				.catch(console.error)
		}
	},
	down: async (queryInterface: QueryInterface, _sequelize: Sequelize) => {
		return queryInterface.bulkDelete('users', { logging: true })
	}
}
