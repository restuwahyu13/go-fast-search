import { QueryInterface, Sequelize, DataTypes } from 'sequelize'

module.exports = {
	up: async (queryInterface: QueryInterface, sequelize: Sequelize) => {
		const tablExist: boolean = await queryInterface.tableExists('users')
		if (!tablExist) {
			await queryInterface.createTable(
				'users',
				{
					id: { type: DataTypes.UUID, primaryKey: true, allowNull: false, unique: true, defaultValue: sequelize.literal('uuid_generate_v4()') },
					name: { type: DataTypes.STRING(200) },
					email: { type: DataTypes.STRING(200) },
					phone: { type: DataTypes.STRING(200) },
					date_of_birth: { type: DataTypes.DATEONLY },
					age: { type: DataTypes.STRING(200) },
					address: { type: DataTypes.TEXT },
					city: { type: DataTypes.STRING(200) },
					state: { type: DataTypes.STRING(200) },
					direction: { type: DataTypes.STRING(200) },
					country: { type: DataTypes.STRING(200) },
					postal_code: { type: DataTypes.STRING(200) },
					is_sync: { type: DataTypes.BOOLEAN, allowNull: false, defaultValue: false },
					created_at: { type: DataTypes.DATE, allowNull: false, defaultValue: sequelize.literal('CURRENT_TIMESTAMP') },
					updated_at: { type: DataTypes.DATE },
					deleted_at: { type: DataTypes.DATE }
				},
				{
					logging: true
				}
			)
		}
	},
	down: async (queryInterface: QueryInterface, _sequelize: Sequelize) => {
		const tableExist: boolean = await queryInterface.tableExists('users')
		if (tableExist) {
			return queryInterface.dropTable('users')
		}
	}
}
