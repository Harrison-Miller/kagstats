import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'kagClassImg',
  pure: true
})
export class KagClassImgPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getImgSource(value);
  }

  getImgSource(kagClass: string) : string {
  
    let classImgMap = {
      knight: 'https://wiki.kag2d.com/images/6/6b/Knight.png',
      archer: 'https://wiki.kag2d.com/images/2/29/Archer.png',
      builder: 'https://wiki.kag2d.com/images/1/14/Builder.png'
    }

    return classImgMap[kagClass];
  }

}
