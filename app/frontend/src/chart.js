// import van from 'vanjs-core'
// import { ChartConfiguration, LibChart } from 'chart.js/auto'

// import { Button } from './buttons.js'
// import { ModalBackground, ModalForeground } from './modal.js';

// const { div, canvas } = van.tags;

// const MakeChartData = () => {
//     return {
//         labels: chartData.val.map(([x, _y]) => x),
//         datasets: [
//           {
//             label: "Global Market Cap ($T)",
//             backgroundColor: "#0eee8c",
//             borderColor: "#0eee8c",
//             data: chartData.val.map(([_x, y]) => y),
//           },
//         ],
//     }
// }

// const MakeChartConfig = (chartData) => {
//     return {
//         type: "line",
//         data: chartData,
//         options: {
//           scales: {
//             x: {
//               title: {
//                 display: true,
//                 text: "Day",
//               },
//             },
//             y: {
//               title: {
//                 display: true,
//                 text: "Market Cap ($T)"
//               },
//               beginAtZero: false,
//               ticks: {
//                 // Include a dollar sign in the ticks
//                 callback: function(value, _index, _ticks) {
//                   // we want value to be rounded to 3 decimal places
//                   // 1.23456789 -> 1.235
//                   return "$" + (value as number).toFixed(3) + "T";
//                 },
//               },
//             },
//           },
//         },
//     }
// }

// export const ChartModal = (config) => {
//     const chart = Chart(config)

//     const closed = van.state(false);
//     const buttonOK = Button({
//         child: "OK",
//         onclick: () => closed.val = true,
//     })

//     return () => closed.val ? null : ModalBackground(
//         ModalForeground({}, chart, buttonOK),
//     )
// }

// const Chart = (config) => {
//     const chart = new LibChart(canvas(), config)

//     return div(chart.canvas)
// }

// export {Chart, ChartModal, MakeChartConfig}