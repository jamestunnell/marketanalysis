import van from "vanjs-core"

const {div, table, td, th, thead, tr} = van.tags

const TableRow = (rowItems) => {
    return tr(
        {class: "border border-solid"},
        rowItems.map(x => td({class: "px-6 py-4"}, x)),
    )
}

const Table = ({columnNames, tableBody}) => {
    const tableHead = thead();

    van.add(tableHead, tr(
        {class: "border-b border-neutral-200 font-medium"},
        columnNames.map(x => th({class: "px-6 py-4"}, x))
    ));

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