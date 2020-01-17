data ="windowSizePerFlowSort.dat"
stats data using 2 name "window_stats"


set xlabel "Packet Count per Flow" 
set ylabel "Cumulative Percentage" 
set tics font "Arial, 18"
set key font "Arial,18"
set key right bottom
set palette gray


c1=0
cumulative_sum1(x)=(c1=c1+x,c1)

set terminal postscript eps color enhanced "Arial" 25
#set term postscript enhanced eps color font ",24" size 5., 7.1
set output "windowSizePerFlow.eps"

set logscale x
#set yr[0:100]
set xrange[1:]




plot data using 1:(cumulative_sum1($2)/window_stats_sum*100)  w l title "window size"  ,\
