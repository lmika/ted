/**
 * The model.
 */
package main


/**
 * An abstract model interface.
 */
type Model interface {
    
    /**
     * The dimensions of the model (width, height).
     */
    GetDimensions()     (int, int)
}
