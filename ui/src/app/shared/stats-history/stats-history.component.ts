import {
  AfterViewInit,
  Component,
  Input,
  OnChanges,
  OnInit,
  SimpleChanges,
  ViewChild,
  ViewChildren
} from '@angular/core';

import {
  ChartComponent,
  ApexAxisChartSeries,
  ApexXAxis,
  ApexLegend
} from 'ng-apexcharts';
import {PlayersService} from '../../services/players.service';
import {MonthlyHittersStats, MonthlyStats} from '../../models';
import {HittersService} from '../../services/hitters.service';
import {first, takeUntil} from "rxjs/operators";

@Component({
  selector: 'app-stats-history',
  templateUrl: './stats-history.component.html',
  styleUrls: ['./stats-history.component.scss', '../../shared/killfeed/killfeed.component.scss']
})
export class StatsHistoryComponent implements OnInit, AfterViewInit {

  constructor(private playersService: PlayersService,
              private hittersService: HittersService) {

  }
  @ViewChildren('chart') chart;
  @ViewChildren('kdrChart') kdrChart;
  @ViewChildren('hittersChart') hittersChart;
  public chartOptions;
  public kdrChartOptions;
  public hittersChartOptions;

  @Input() playerID: number;

  loading = 2;

  stats: MonthlyStats[];
  hitterStats: MonthlyHittersStats[];
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

  hitterInfo = {
    // distinct hitters

    // knight
    sword: {
      color: '#FF9800FF',
      rank: 0,
    },
    bomb: {
      color: '#FF9800C0',
      rank: 3,
    },
    shield: {
      color: '#FF980080',
      rank: 27,
    },

    // archer
    arrow: {
      color: '#2E93fAFF',
      rank: 1,
    },
    bomb_arrow: {
      color: '#2E93fAC0',
      rank: 10,
    },

    // builder
    spikes: {
      color: '#66DA26FF',
      rank: 6,
    },
    pickaxe: {
      color: '#66DA26C0',
      rank: 9,
    },
    drill: {
      color: '#66DA2680',
      rank: 11,
    },
    saw: {
      color: '#66DA2640',
      rank: 15,
    },

    // misc
    stomp: {
      color: '#6E7F80FF',
      rank: 5,
    },
    fall: {
      color: '#6E7F80C0',
      rank: 7,
    },
    keg: {
      color: '#9370DBFF',
      rank: 4,
    },
    mine: {
      color: '#9370DBC0',
      rank: 14,
    },
    bite: {
      color: '#7FFFD4FF',
      rank: 20,
    },
    fire: {
      color: '#DC143CFF',
      rank: 17,
    },
    catapult_stones: {
      color: '#CB6843FF',
      rank: 16,
    },
    ballista_bolt: {
      color: '#CB6843C0',
      rank: 22,
    },
    flying: {
      color: '#CB684380',
      rank: 25,
    },
    boulder: {
      color: '#A0522DFF',
      rank: 23,
    },
    drowning: {
      color: '#20B2AAFF',
      rank: 18,
    },
    water: {
      color: '#20B2AAC0',
      rank: 31,
    },
    sudden_gib: {
      color: '#9370DB80',
      rank: 30,
    },

    // non-distinct hitters
    died: {
      color: '#2e222f',
      rank: 12,
    },
    crushing: {
      color: '#453347',
      rank: 13,
    },
    suicide: {
      color: '#2e222f',
      rank: 19,
    },
    water_stun: {
      color: '#4d65b4',
      rank: 29,
    },
    water_stun_force: {
      color: '#4d65b4',
      rank: 26,
    },
    burn: {
      color: '#b33831',
      rank: 17,
    },
    explosion: {
      color: '#e968dc',
      rank: 24,
    },
    mine_special: {
      color: '#b553ab',
      rank: 8,
    },
    catapult_boulder: {
      color: '#434341',
      rank: 28,
    },

    // unused
    ram: {
      color: '#f9c22b',
      rank: 40
    },
    stab: {
      color: '#2E93fA',
      rank: 41,
    },
    muscles: {
      color: '#f9c22b',
      rank: 42,
    },
  };

  ngAfterViewInit() {
    this.chart.changes.pipe(first()).subscribe(result => {
      const chart = result.first;
      setTimeout(() => {
        chart.toggleSeries('Archer Deaths');
        chart.toggleSeries('Builder Deaths');
        chart.toggleSeries('Knight Deaths');
      }, 50);
    });

    this.hittersChart.changes.pipe(first()).subscribe( result => {
      const chart = result.first;
      setTimeout(() => {
        chart.toggleSeries('Drowning');
        chart.toggleSeries('Shark');
        chart.toggleSeries('Splatter');
        chart.toggleSeries('Water');
        chart.toggleSeries('Shield');
        chart.toggleSeries('Scroll Of Carnage');
      }, 50);
    });
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

    this.hittersChartOptions = {
      series: [
      ],
      chart: {
        height: 800,
        type: 'bar',
        stacked: true,
        stackType: '100%',
        animations: {
          enabled: false,
        },
        toolbar: {
          show: false,
        },
      },
      title: {
        text: 'Monthly Weapons'
      },
      plotOptions: {
        bar: {
          columnWidth: '90%',
        },
      },
      xaxis: {
        type: 'category',
        categories: [],
        labels: {
          rotate: -45,
          rotateAlways: true,
        },
      },
      yaxis: {
        show: false,
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

    this.hittersService.getMonthlyHitters(this.playerID).subscribe( hitterStats => {
      this.hitterStats = hitterStats;
      this.hitterStats.sort((a, b) => {
        if (a.year === b.year) {
          return a.month - b.month;
        }
        return a.year - b.year;
      });
      this.hitterStats.slice(0, 11);
      this.updateHitterStats();
    });
  }

  categoryName(month: number, year: number): string {
    switch (month) {
      case 1:
        return 'Jan ' + year;
      case 2:
        return 'Feb ' + year;
      case 3:
        return 'Mar ' + year;
      case 4:
        return 'Apr ' + year;
      case 5:
        return 'May ' + year;
      case 6:
        return 'Jun ' + year;
      case 7:
        return 'Jul ' + year;
      case 8:
        return 'Aug ' + year;
      case 9:
        return 'Sep ' + year;
      case 10:
        return 'Oct ' + year;
      case 11:
        return 'Nov ' + year;
      case 12:
        return 'Dec ' + year;
    }
    return '' + year;
  }

  getKDR(kills: number, deaths: number): string {
    return (kills / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }

  updateStats() {
    for (const stat of this.stats) {
      const category = this.categoryName(stat.month, stat.year);
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

    this.chartOptions.series = [this.knightKills,
      this.knightDeaths,
      this.archerKills,
      this.archerDeaths,
      this.builderKills,
      this.builderDeaths];

    /*
    this.chart.appendSeries(this.knightKills);
    this.chart.appendSeries(this.knightDeaths);
    this.chart.appendSeries(this.archerKills);
    this.chart.appendSeries(this.archerDeaths);
    this.chart.appendSeries(this.builderKills);
    this.chart.appendSeries(this.builderDeaths);
     */

    this.kdrChartOptions.series = [this.knightKDR, this.archerKDR, this.builderKDR];

    /*
    this.kdrChart.appendSeries(this.knightKDR);
    this.kdrChart.appendSeries(this.archerKDR);
    this.kdrChart.appendSeries(this.builderKDR);
     */

    this.loading -= 1;
  }

  hitterDescription(name: string): string {
    switch (name) {
      case 'catapult_stones':
        return 'Catapult';
      case 'bite':
        return 'Shark';
      case 'flying':
        return 'Splatter';
      case 'ballista_bolt':
        return 'Ballista';
      case 'sudden_gib':
        return 'Scroll Of Carnage';
      case 'fall':
        return 'Pushed Off A Cliff';
    }

    let description = '';
    const parts = name.split('_');
    for (const i in parts) {
      if (i !== '0') {
        description += ' ';
      }
      description += parts[i].charAt(0).toUpperCase() + parts[i].slice(1);
    }
    return description;
  }

  combineHitter(name: string): string {
    switch (name) {
      case 'died':
      case 'suicide':
      case 'crushing':
        return 'fall';
      case 'fire':
      case 'burn':
        return 'fire';
      case 'explosion':
        return 'bomb';
      case 'mine_special':
        return 'mine';
      case 'catapult_boulder':
      case 'catapult_stones':
        return 'catapult_stones';
      case 'water':
      case 'water_stun':
      case 'water_stun_force':
        return 'water';
    }
    return name;
  }

  updateHitterStats() {
    const series = {};
    for (const stat of this.hitterStats) {
      this.hittersChartOptions.xaxis.categories.push(this.categoryName(stat.month, stat.year));
    }

    for (const stat of this.hitterStats) {
      // create counts list combining duplicate hitter types
      const counts = {};
      for (const key in stat) {
        if (key !== 'player' && key !== 'month' && key !== 'year' && key !== 'ram' && key !== 'muscles' && key != 'stab') {
          const hitterName = this.combineHitter(key);
          // if hitterName doesn't exist create it
          if (!(hitterName in counts)) {
            // if 0 null
            if (stat[key] === 0) {
              counts[hitterName] = null;
            } else {
              counts[hitterName] = stat[key];
            }
          } else {
            // if hitterName in counts
            // if previously set to null override it
            if (counts[hitterName] == null) {
              counts[hitterName] = 0;
            }
            counts[hitterName] += stat[key];
          }
        }
      }

      // use counts list to populate series
      for (const key in counts) {
        if (!(key in series)) {
          series[key] = {
            key,
            name: this.hitterDescription(key),
            data: [],
            color: this.hitterInfo[key].color,
          };
        }
        series[key].data.push(counts[key]);
      }
    }

    const sortedSeries = [];
    for (const key in series) {
      sortedSeries.push(series[key]);
    }

    sortedSeries.sort((a, b) => {
      const rankA = this.hitterInfo[a.key].rank;
      const rankB = this.hitterInfo[b.key].rank;
      return rankA - rankB;
    });

    for (const s of sortedSeries) {
      this.hittersChartOptions.series.push(s);
      // this.hittersChart.appendSeries(s);
    }

    this.loading -= 1;
  }

}
