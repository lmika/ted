(defn gotoCell [r c]
    (cellRow r)
    (cellCol c))

(defn cellValueAt [r c v]
    (gotoCell r c)
    (cellValue v))


(modelResize 10 10)
(for [(def r 0) (< r 10) (def r (+ r 1))]
    (for [(def c 0) (< c 10) (def c (+ c 1))]
        (cellValueAt r c (concat (str r) "," (str c)))))