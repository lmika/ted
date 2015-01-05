// Various utilities

package ui


// Returns the maximum value of either x or y.
func intMax(x, y int) int {
    if (x < y) {
        return y
    } else {
        return x
    }
}

// Returns the minimum value of either x or y.
func intMin(x, y int) int {
    if (x > y) {
        return y
    } else {
        return x
    }
}

// Returns the value capped between two limits.
func intMinMax(x, min, max int) int {
    if (x < min) {
        return min
    } else if (x > max) {
        return max
    } else {
        return x
    }
}