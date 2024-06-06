import van from "vanjs-core"

import { Route } from 'vanjs-router'
import GraphsPage from './src/graphs.js'
import GraphPage from './src/graph.js'
import NavBar from './src/navbar.js'
import GraphSettings from "./src/graphsettings.js"

import './index.css';

const RouteHome = () => {
    return Route(
        { name: 'home' },
        NavBar({currentRoute: 'home'}),
    )
}

const RouteGraphs = () => {
    const graphID = van.state('')
    const settings = new GraphSettings()
    
    return Route(
        {
            name: 'graphs',
            onFirst() {
            },
            onLoad(route) {
                settings.load()
                
                if (route.args.length == 0) {
                    graphID.val = '';

                    return
                }

                graphID.val = route.args[0]
            }
        },
        () => {
            if (graphID.val === '') {
                return NavBar({currentRoute: 'graphs'})
            }

            return NavBar({currentRoute: `graphs/${graphID.val}`})
        },
        settings.render(),
        () => {
            if (graphID.val === '') {
                return GraphsPage()
            }

            return GraphPage({id: graphID.val, settings})
        }
    )
}

van.add(document.body, RouteHome(), RouteGraphs())
