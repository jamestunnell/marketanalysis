import van from "vanjs-core"

import {Get} from './backend.js'

const {button, div, p} = van.tags

const getGraphs = async () => {
    console.log("getting graphs");

    const resp = await Get('/graphs');

    if (resp.status != 200) {
        console.log("failed to get graphs", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log("received %d graphs", d.graphs.length);

    return d.graphs;
}

// const delGraph = async (symbol) => {
//     console.log("deleting graph");

//     const resp = await fetch(`${BASE_URL}/graphs/${symbol}`, {
//         method: 'DELETE',
//         credentials: 'same-origin'
//     });

//     console.log('delete graph result:', resp.status)

//     return resp.status === 204 
// }

const GraphSidebarItem = ({data, state}) => {
    const deleted = van.state(false);
    const itemState = {
        id: van.state(data.id),
    }
    const itemClass = van.derive(() => {
        const isSelected = state.selectedID.val == itemState.id.val;

        return `block divide-y md:px-4 md:py-2 ${isSelected ? "text-indigo-500" : "text-gray-500"}`
    });

    return () => deleted.val ? null : button(
        {
            class: itemClass,
        },
        div(
            p({class: "font-semibold"}, data.name),
            p({class: "italic"}, truncateID(data.id))
        )
    );
}

const ID_PREVIEW_LEN = 8;

const truncateID = (id) => {
    if (id.length > ID_PREVIEW_LEN) {
        return id.substring(0, ID_PREVIEW_LEN) + "..."
    }
    
    return id
}

const Graphs = () => {
    const state = {
        selectedID: van.state(""),
    }

    const sidebarArea = div(
        {class:"flex flex-col flex-nowrap overflow-y-scroll"},
    )

    const mainArea = div(
        {class: "h-screen flex flex-col px-6 py-4"},
    )
    
    getGraphs().then(
        (items) => {
            const sidebarItems = items.map(data => GraphSidebarItem({data: data, state: state}));

            van.add(sidebarArea, sidebarItems);
            // van.add(sidebarArea, AddNewButton({sidebar: sidebarArea, state: state}));
        }
    );

    return div(
        {class: "h-screen flex"},
        sidebarArea,
        mainArea,
    );
}

export default Graphs