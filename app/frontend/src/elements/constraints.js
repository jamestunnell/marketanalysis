
const constraintMinMax = ({constraint, step, onMin, onMax}) => {
    switch (constraint.type) {
        case 'less':
            onMax(constraint.limits[0] - step)
            break
        case 'lessEqual':
            onMax(constraint.limits[0])
            break
        case 'greater':
            onMin(constraint.limits[0] + step)
            break
        case 'greaterEqual':
            onMin(constraint.limits[0])
            break
        case 'rangeIncl':
            onMin(constraint.limits[0])
            onMax(constraint.limits[1])
            break
        case 'rangeExcl':
            onMin(constraint.limits[0])
            onMax(constraint.limits[1] - step)
            break
    }
}

export {constraintMinMax}