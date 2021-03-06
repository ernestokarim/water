
(println " > Common operators")
(+ 4 3)
(- 3 2)
(* 2 3)
(/ 10 2)
(% 11 5)
(- 5)
(- -5)

(println "")
(println " > Greater than")
(> 3 2)
(> 2 3)

(println "")
(println " > Less than")
(< 2 3)
(< 3 2)

(println "")
(println " > Greater than or equal to")
(>= 3 2)
(>= 3 3)
(>= 2 3)

(println "")
(println " > Less than or equal to")
(<= 2 3)
(<= 3 3)
(<= 3 2)

(println "")
(println " > Equal to")
(= 3 3)
(= 2 3)

(println "")
(println " > Not")
(not #t)
(not #f)
(not (not #t))

###########################################################

 > Common operators
7
1
6
5
1
-5
5

 > Greater than
true
false

 > Less than
true
false

 > Greater than or equal to
true
true
false

 > Less than or equal to
true
true
false

 > Equal to
true
false

 > Not
false
true
true
