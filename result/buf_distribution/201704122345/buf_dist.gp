buf_dist1data = "buf_dist1.dat"
stats buf_dist1data using 2 name "buf_dist1stats"
buf_dist01data = "buf_dist01.dat"
stats buf_dist01data using 2 name "buf_dist01stats"
buf_dist001data = "buf_dist001.dat"
stats buf_dist001data using 2 name "buf_dist001stats"
buf_dist0001data = "buf_dist0001.dat"
stats buf_dist0001data using 2 name "buf_dist0001stats"
buf_dist00001data = "buf_dist00001.dat"
stats buf_dist00001data using 2 name "buf_dist00001stats"
buf_dist000001data = "buf_dist000001.dat"
stats buf_dist000001data using 2 name "buf_dist000001stats"
buf_dist0000001data = "buf_dist0000001.dat"
stats buf_dist0000001data using 2 name "buf_dist0000001stats"

set xlabel "Packets of Buffer"
set ylabel "Cumulative Percentage"
set tics font "Arial, 10"

c1=0
cumulative_sum1(x)=(c1=c1+x,c1)
c01=0
cumulative_sum01(x)=(c01=c01+x,c01)
c001=0
cumulative_sum001(x)=(c001=c001+x,c001)
c0001=0
cumulative_sum0001(x)=(c0001=c0001+x,c0001)
c00001=0
cumulative_sum00001(x)=(c00001=c00001+x,c00001)
c000001=0
cumulative_sum000001(x)=(c000001=c000001+x,c000001)
c0000001=0
cumulative_sum0000001(x)=(c0000001=c0000001+x,c0000001)

set logscale x
#set yr[0:100]
set xrange[1:]

plot buf_dist1data using 1:(cumulative_sum1($2)/buf_dist1stats_sum*100)  w l title "time width = 1.0",\
     buf_dist01data using 1:(cumulative_sum01($2)/buf_dist01stats_sum*100)  w l title "time width = 0.1",\
     buf_dist001data using 1:(cumulative_sum001($2)/buf_dist001stats_sum*100)  w l title "time width = 0.01",\
     buf_dist0001data using 1:(cumulative_sum0001($2)/buf_dist0001stats_sum*100)  w l title "time width = 0.001",\
     buf_dist00001data using 1:(cumulative_sum00001($2)/buf_dist00001stats_sum*100)  w l title "time width = 0.0001",\
     buf_dist000001data using 1:(cumulative_sum000001($2)/buf_dist000001stats_sum*100)  w l title "time width = 0.00001",\
     buf_dist0000001data using 1:(cumulative_sum0000001($2)/buf_dist0000001stats_sum*100)  w l title "time width = 0.00001",\
 