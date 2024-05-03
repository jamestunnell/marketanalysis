import DataTable from 'react-data-table-component';

const columns = [
	{
		name: 'Symbol',
		selector: row => row.symbol,
    sortable: true,
	},
	{
		name: 'Time Zone',
		selector: row => row.timeZone,
    sortable: true,
	},
	{
		name: 'Open',
		selector: row => row.open,
	},
	{
		name: 'Close',
		selector: row => row.close,
	},
];

export default function SecuritiesTable({securities}) {
  return (
		<DataTable
			columns={columns}
			data={securities}
      pagination
		/>
	);
}
