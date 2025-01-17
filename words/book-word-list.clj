(defn run [[book-uri output-file]]
  (->> (slurp book-uri)
       (re-seq #"[\w|’|']+")
       (map clojure.string/lower-case)
       (filter #(and (> 7 (count %)) (< 1 (count %))))
       (frequencies)
       (sort-by val)
       (take-last 500)
       (map key)
       (shuffle)
       (clojure.string/join "\n")
       (#(clojure.string/replace % #"’" "'"))
       (#(clojure.string/replace % #"_" ""))
       (spit output-file)))

(run *command-line-args*)

(println "Done!")
