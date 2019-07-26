import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'unkownAvatar',
  pure: true
})
export class UnkownAvatarPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getPortrait(value);
  }

  getPortrait(userID: number): string {
    let portraitNum = userID%3 + 1;
    return `/assets/portrait${portraitNum}.png`
  }

}
