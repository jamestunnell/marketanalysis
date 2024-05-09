import van from "vanjs-core"

import {Get} from './backend.js'

const {a, div, li, ul} = van.tags

const getGraphs = async () => {
    console.log("getting graphs");

    const resp = await Get('/graphs');

    if (!resp.ok) {
        console.log("non-ok get graphs response %d: %s", resp.status, resp.text());

        return []
    }
    
    const d = await resp.json()

    console.log("get graphs response JSON: %o", d)

    return d.graphs;
}

const delGraph = async (symbol) => {
    console.log("deleting graph");

    const resp = await fetch(`${BASE_URL}/graphs/${symbol}`, {
        method: 'DELETE',
        credentials: 'same-origin'
    });

    console.log('delete graph result:', resp.status)

    return resp.status === 204 
}

const ListItem = ({id}) => {
    const deleted = van.state(false)
    return () => deleted.val ? null : li(
        div(
            {class: "flex flex-row gap-4"},
            id,
            a(
                {
                    onclick: () => {
                        if (delGraph(id)) {
                            deleted.val = true
                        }
                    }
                },
                "❌",
            ),
        )
    )
}

const Graphs = () => {
    const listDom = ul({class:"p-4"})

    getGraphs().then(
        (items) => {
            console.log("found %d graphs", items.length);
            
            const listItems = items.map(x => ListItem({id: x.id}))
            
            van.add(listDom, listItems)
        }
    ).catch(error => {
        console.log("failed to get graphs: " + error);
    });

    const dom = div();

    return div(
        {class: "h-screen w-screen p-4"},
        dom,
        listDom
    );
}

export default Graphs