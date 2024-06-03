import van from "vanjs-core"
import { routeTo } from 'vanjs-router'
import {Modal} from "vanjs-ui"
import { v4 as uuidv4 } from 'uuid';
import { uniqueNamesGenerator, adjectives, animals } from 'unique-names-generator';

import {Delete, Get, PostJSON} from './backend.js'
import { Button, ButtonIcon, ButtonCancel } from "./buttons.js";
import { ButtonGroup } from './buttongroup.js'
import { IconAdd, IconDelete, IconView } from './icons.js'
import {Table, TableRow} from './table.js'
import truncateString from "./truncatestring.js";

const {div, input, label, p, tbody} = van.tags

const getGraphs = async () => {
    console.log("getting graphs");

    const resp = await Get('/graphs');

    if (resp.status != 200) {
        console.log("failed to get graphs", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log(`received ${d.graphs.length} graphs`, d.graphs);

    return d.graphs;
}

const createGraph = async (item) => {
    console.log("creating graph", item);

    const resp = await PostJSON({route: '/graphs', object: item});

    if (resp.status != 204) {
        console.log("failed to create graph", await resp.json());

        return false
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log(`created graph %s`, item.id);

    return true;
}

const deleteGraph = async (id) => {
    console.log("deleting graph %s", id);

    const resp = await Delete(`/graphs/${id}`);

    if (resp.status != 204) {
        console.log("failed to delete graph", await resp.json());

        return false
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log(`deleted graph %s`, id);

    return true;
}

const RandomName = () => {
    return uniqueNamesGenerator({ dictionaries: [adjectives, animals] });
}

const GraphNameForm = ({onOK, onCancel}) => {
    const name = van.state(RandomName())

    return div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Name"),
        label({for: "name"}, "Name"),
        input({
            id: "name",
            class: "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring",
            type: "text",
            value: name,
            placeholder: "Unique, non-empty name",
            oninput: e => name.val = e.target.value,
        }),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({child: "Cancel", onclick: onCancel}),
            Button({child: "OK", onclick: ()=> onOK({name: name.val})}),
        ),
    )
}

const GraphTableRow = ({id, name}) => {
    const deleted = van.state(false)

    const viewBtn = ButtonIcon({
        icon: IconView(),
        // text: "View",
        onclick: () => routeTo('graphs', [id]),
    });
    const deleteBtn = ButtonIcon({
        icon: IconDelete(),
        // text: "Delete",
        onclick: () => {
            deleteGraph(id).then(ok => {
                if (ok) {
                    deleted.val = true
                }
            })
        },
    });

    const buttons = ButtonGroup({buttons: [viewBtn, deleteBtn]});
    const rowItems = [name, truncateString(id, 8), buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const GraphsPage = () => {
    const columnNames = ["Name", "ID", ""]
    const tableBody = tbody({class:"table-auto"});

    getGraphs().then((graphs) => {
        const rows = graphs.map(g => GraphTableRow({id: g.id, name: g.name}));

        van.add(tableBody, rows);
    });

    const addIcon = IconAdd()
    
    addIcon.classList.add("text-xl")

    const newGraphBtn = ButtonIcon({
        icon: addIcon,
        onclick: () => {
            const closed = van.state(false)

            van.add(
                document.body,
                Modal({closed},
                    GraphNameForm({
                        onOK: ({name})=> {
                            const id  = uuidv4();
                            const graphItem = {id: id, name: name, blocks: [], connections: []};

                            createGraph(graphItem).then((ok) => {
                                if (ok) {
                                    van.add(tableBody, GraphTableRow({id: id, name: name, }));
                                    
                                    closed.val = true;
                                }
                            });
                        },
                        onCancel: () => {
                            closed.val = true
                        }
                    }),
                ),
            );
        },
    });

    return div(
        {class: "container p-4 w-full flex flex-col divide-y divide-gray-400"},
        div(
            {class: "flex flex-col mt-4"},
            div(
                {class: "grid grid-cols-2"},
                div(
                    {class: "flex flex-row p-2"},
                    p({class: "p-3 m-1 text-xl font-medium"}, "Graphs"),
                ),
                div(
                    {class: "flex flex-row-reverse p-2"},
                    newGraphBtn,
                )
            ),
            Table({columnNames: columnNames, tableBody: tableBody}),
        )
    )
}

export default GraphsPage;