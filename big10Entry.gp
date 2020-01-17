set xlabel "Time [s]"
set ylabel "Number of Entrys"
set tics font "Arial, 10"
set title font"Arial,5"
set label font "Arial,5"
set key font "Arial,10"

#set logscale y
set yr[0:1400]
set xrange[1:]
set terminal pngcairo
set output "big10EntrySinet.png"

plot "big10Entrysinet.dat" using 1:2  w l title "1st",\
     "big10Entrysinet.dat" using 1:3  w l title "2nd",\
     "big10Entrysinet.dat" using 1:4  w l title "3rd",\
     "big10Entrysinet.dat" using 1:5  w l title "4th",\
     "big10Entrysinet.dat" using 1:6  w l title "5th",\
     "big10Entrysinet.dat" using 1:7  w l title "6th",\
     "big10Entrysinet.dat" using 1:8  w l title "7th",\
     "big10Entrysinet.dat" using 1:9  w l title "8th",\
     "big10Entrysinet.dat" using 1:10  w l title "9th",\
     "big10Entrysinet.dat" using 1:11 w l title "10th"