import van from "vanjs-core"

const {div, table, tbody, td, th, thead, tr} = van.tags

const TableRow = (rowItems) => {
    const deleted = van.state(false)
    return () => deleted.val ? null : tr(
        {class: "border border-solid"},
        rowItems.map(x => td({class: "px-6 py-4"}, x)),
    )
}

const Table = ({columnNames, rows}) => {
    const tableHead = thead();
    const tableBody = tbody({class:"table-auto"});

    van.add(tableHead, tr(
        {class: "border-b border-neutral-200 font-medium"},
        columnNames.map(x => th({class: "px-6 py-4"}, x))
    ));
    van.add(tableBody, rows);

    return div(
        {class: "w-full p-4 space-y-6"},
        table(
            {class:"table-auto min-w-full text-left text-sm"},
            tableHead,
            tableBody,
        ),
    );
}

export {TableRow, Table};