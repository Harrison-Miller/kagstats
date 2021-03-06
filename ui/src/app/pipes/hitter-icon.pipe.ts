import { Pipe, PipeTransform } from '@angular/core';
import { Hitters } from  '../hitters';
import { HttpClientJsonpModule } from '@angular/common/http';

@Pipe({
  name: 'hitterIcon',
  pure: true
})
export class HitterIconPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getHitterIcon(value);
  }

  getHitterIcon(hitter: number): string {
    // logic mostly taken from Base/Rules/CommonScripts/KillMessages.as
    switch(hitter) {
      case Hitters.nothing:
      case Hitters.suicide:
      case Hitters.crush:
        return Hitters[Hitters.fall];
      case Hitters.fire:
      case Hitters.burn:
        return Hitters[Hitters.fire];
      case Hitters.explosion:
        return Hitters[Hitters.bomb];
      case Hitters.mine_special:
        return Hitters[Hitters.mine];
      case Hitters.cata_boulder:
      case Hitters.cata_stones:
        return Hitters[Hitters.cata_stones];
      case Hitters.water:
      case Hitters.water_stun:
      case Hitters.water_stun_force:
        return Hitters[Hitters.water];
    }

    return Hitters[hitter];
  }

}
