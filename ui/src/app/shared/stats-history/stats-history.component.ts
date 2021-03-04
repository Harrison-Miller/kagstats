import {Component, Input, OnInit, ViewChild} from '@angular/core';

import {
  ChartComponent,
} from 'ng-apexcharts';
import {PlayersService} from '../../services/players.service';
import {MonthlyStats} from '../../models';

@Component({
  selector: 'app-stats-history',
  templateUrl: './stats-history.component.html',
  styleUrls: ['./stats-history.component.scss']
})
export class StatsHistoryComponent implements OnInit {
  public chartOptions;
  public kdrChartOptions;

  @Input() playerID: number;

  stats: MonthlyStats[];
  categories = [];
  archerKills = {
    name: 'Archer Kills',
    color: '#2E93fA',
    data: [],
  };
  archerDeaths = {
    name: 'Archer Deaths',
    color: '#1b5896',
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
    data: []
  };
  knightKills = {
    name: 'Knight Kills',
    color: '#FF9800',
    data: [],
  };
  knightDeaths = {
    name: 'Knight Deaths',
    color: '#995b00',
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
        this.archerKills,
        this.archerDeaths,
        this.builderKills,
        this.builderDeaths,
        this.knightKills,
        this.knightDeaths
      ],
      chart: {
        height: 400,
        type: 'area',
        stacked: true,
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
      dataLabels: {
        enabled: false,
      },
      title: {
        text: 'Monthly Kills and Deaths'
      },
      xaxis: {
        categories: this.categories,
      },
    };

    this.kdrChartOptions = {
      series: [
        this.archerKDR,
        this.builderKDR,
        this.knightKDR
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
      dataLabels: {
        enabled: true,
      },
      markers: {
        size: 8,
        shape: 'circle',
      },
      title: {
        text: 'Monthly KDR'
      },
      xaxis: {
        categories: this.categories,
      },
      yaxis: {
        logarithmic: true,
      }
    };


    this.playersService.getMonthlyStatus(this.playerID).subscribe( monthlyStats => {
      this.stats = monthlyStats;
      this.stats.sort((a, b) => {
        if (a.year === b.year) {
          return b.month - a.month;
        }
        return b.year - a.year;
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
      this.categories.push(this.categoryName(stat));
      this.archerKills.data.push(stat.archerKills);
      this.archerDeaths.data.push(stat.archerDeaths);
      this.archerKDR.data.push(this.getKDR(stat.archerKills, stat.archerDeaths));
      this.builderKills.data.push(stat.builderKills);
      this.builderDeaths.data.push(stat.builderDeaths);
      this.builderKDR.data.push(this.getKDR(stat.builderKills, stat.builderDeaths));
      this.knightKills.data.push(stat.knightKills);
      this.knightDeaths.data.push(stat.knightDeaths);
      this.knightKDR.data.push(this.getKDR(stat.knightKills, stat.knightDeaths));

    }

    this.chartOptions.series = [this.archerKills, this.archerDeaths, this.builderKills, this.builderDeaths, this.knightKills, this.knightDeaths];
    this.chartOptions.xaxis = { categories: this.categories };

    this.kdrChartOptions.series = [this.archerKDR, this.builderKDR, this.knightKDR];
    this.kdrChartOptions.xaxis = { categories: this.categories };
  }

}
