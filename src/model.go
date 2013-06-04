/**
 * The model.
 */
package main


/**
 * An abstract model interface.
 */
type IModel interface (

    /**
     * The dimensions of the model (width, height).
     */
    GetDimensions()     (int, int)
)
