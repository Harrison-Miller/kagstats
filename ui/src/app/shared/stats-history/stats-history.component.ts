import {Component, Input, OnInit, ViewChild} from '@angular/core';

import {
  ChartComponent,
  ApexAxisChartSeries,
  ApexXAxis
} from 'ng-apexcharts';
import {PlayersService} from '../../services/players.service';
import {MonthlyStats} from '../../models';

@Component({
  selector: 'app-stats-history',
  templateUrl: './stats-history.component.html',
  styleUrls: ['./stats-history.component.scss']
})
export class StatsHistoryComponent implements OnInit {
  @ViewChild('chart') chart;
  @ViewChild('kdrChart') kdrChart;
  public chartOptions;
  public kdrChartOptions;

  @Input() playerID: number;

  stats: MonthlyStats[];
  archerKills = {
    name: 'Archer Kills',
    color: '#2E93fA',
    type: 'area',
    data: [],
  };
  archerDeaths = {
    name: 'Archer Deaths',
    color: '#1b5896',
    type: 'line',
    enabled: false,
    data: [],
  };
  archerKDR = {
    name: 'Archer KDR',
    color: '#2E93fA',
    data: []
  };
  builderKills = {
    name: 'Builder Kills',
    color: '#66DA26',
    type: 'area',
    data: []
  };
  builderDeaths = {
    name: 'Builder Deaths',
    color: '#3d8216',
    data: []
  };
  builderKDR = {
    name: 'Builder KDR',
    color: '#66DA26',
    type: 'line',
    data: []
  };
  knightKills = {
    name: 'Knight Kills',
    color: '#FF9800',
    type: 'area',
    data: [],
  };
  knightDeaths = {
    name: 'Knight Deaths',
    color: '#995b00',
    type: 'line',
    data: [],
  };
  knightKDR = {
    name: 'Knight KDR',
    color: '#FF9800',
    data: []
  };

  constructor(private playersService: PlayersService) {

  }

  ngOnInit() {
    this.chartOptions = {
      series: [
      ],
      chart: {
        height: 400,
        type: 'line',
        animations: {
          enabled: false,
        },
        selection: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
        zoom: {
          enabled: false,
        }
      },
      stroke: {
        dashArray: [0, 10, 0, 10, 0, 10]
      },
      fill: {
        opacity: [0.25, 1, 0.25, 1, 0.25],
      },
      dataLabels: {
        enabled: true,
      },
      title: {
        text: 'Monthly Kills and Deaths'
      },
      xaxis: {
        type: 'category',
        labels: {
          rotate: -45,
          rotateAlways: true,
        },
      },
    };

    this.kdrChartOptions = {
      series: [
      ],
      chart: {
        height: 400,
        type: 'line',
        animations: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      markers: {
        size: 8,
        shape: 'circle',
      },
      title: {
        text: 'Monthly KDR'
      },
      xaxis: {
        type: 'category',
        labels: {
          rotate: -45,
          rotateAlways: true,
        },
      },
      yaxis: {
        min: 0,
        max: 5,
      }
    };


    this.playersService.getMonthlyStatus(this.playerID).subscribe(monthlyStats => {
      this.stats = monthlyStats;
      this.stats.sort((a, b) => {
        if (a.year === b.year) {
          return a.month - b.month;
        }
        return a.year - b.year;
      });
      this.stats.slice(0, 11);
      this.updateStats();
    });
  }

  categoryName(stat: MonthlyStats): string {
    switch (stat.month) {
      case 1:
        return 'Jan ' + stat.year;
      case 2:
        return 'Feb ' + stat.year;
      case 3:
        return 'Mar ' + stat.year;
      case 4:
        return 'Apr ' + stat.year;
      case 5:
        return 'May ' + stat.year;
      case 6:
        return 'Jun ' + stat.year;
      case 7:
        return 'Jul ' + stat.year;
      case 8:
        return 'Aug ' + stat.year;
      case 9:
        return 'Sep ' + stat.year;
      case 10:
        return 'Oct ' + stat.year;
      case 11:
        return 'Nov ' + stat.year;
      case 12:
        return 'Dec ' + stat.year;
    }
    return '' + stat.year;
  }

  getKDR(kills: number, deaths: number): string {
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  updateStats() {
    for (const stat of this.stats) {
      const category = this.categoryName(stat);
      this.archerKills.data.push({x: category, y: stat.archerKills});
      this.archerDeaths.data.push({x: category, y: stat.archerDeaths});
      this.archerKDR.data.push({x: category, y: this.getKDR(stat.archerKills, stat.archerDeaths)});
      this.builderKills.data.push({x: category, y: stat.builderKills});
      this.builderDeaths.data.push({x: category, y: stat.builderDeaths});
      this.builderKDR.data.push({x: category, y: this.getKDR(stat.builderKills, stat.builderDeaths)});
      this.knightKills.data.push({x: category, y: stat.knightKills});
      this.knightDeaths.data.push({x: category, y: stat.knightDeaths});
      this.knightKDR.data.push({x: category, y: this.getKDR(stat.knightKills, stat.knightDeaths)});
    }

    this.chart.appendSeries(this.knightKills);
    this.chart.appendSeries(this.knightDeaths);
    this.chart.appendSeries(this.archerKills);
    this.chart.appendSeries(this.archerDeaths);
    this.chart.appendSeries(this.builderKills);
    this.chart.appendSeries(this.builderDeaths);


    this.chart.toggleSeries('Archer Deaths');
    this.chart.toggleSeries('Builder Deaths');
    this.chart.toggleSeries('Knight Deaths');

    this.kdrChart.appendSeries(this.knightKDR);
    this.kdrChart.appendSeries(this.archerKDR);
    this.kdrChart.appendSeries(this.builderKDR);

  }
}
