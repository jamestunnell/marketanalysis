import van from "vanjs-core"

import { Route } from 'vanjs-router'
import GraphsPage from './src/graphs.js'
import GraphPage from './src/graph.js'
import NavBar from './src/navbar.js'

import './index.css';

const {div} = van.tags

const RouteHome = () => {
    return Route(
        { name: 'home' },
        NavBar({currentRoute: 'home'}),
    )
}

const RouteGraphs = () => {
    const page = div({class:"overflow-hidden"})

    return Route(
        {
            name: 'graphs',
            onFirst() {
            },
            onLoad(route) {
                console.log("loaded route", route)

                while (page.firstChild) {
                    page.removeChild(page.firstChild)
                }

                if (route.args.length === 0) {
                    van.add(page, NavBar({currentRoute: 'graphs'}))
                    van.add(page, GraphsPage())
                } else {
                    van.add(page, NavBar({currentRoute: `graphs/${route.args[0]}`}))
                    van.add(page, GraphPage(route.args[0]))
                }
            }
        },
        page
    )
}

van.add(document.body, RouteHome(), RouteGraphs())
