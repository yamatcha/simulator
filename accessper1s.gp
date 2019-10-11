set xlabel "Time [s]"
set ylabel "TCAM access count"
set tics font "Arial, 10"
set title font"Arial,5"
set label font "Arial,5"
set key font "Arial,10"

#set logscale y
#set yr[0:100]
set xrange[1:]

plot "test0.txt" using 1:2  w l title "buffer=580",\
    "test1.txt" using 1:2  w l title "buffer=1000",\
    "accessper1s_001.dat" using 1:2  w l title "time\_width=0.01",\
    "accessper1s_0.dat" using 1:2  w l title "time\_width=0",\
