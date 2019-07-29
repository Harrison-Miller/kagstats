import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'kagRoleColor',
  pure: true
})
export class KagRoleColorPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getKagRoleColor(value);
  }

  getKagRoleColor(role: number): string {
    if(role==1 || role==4) {
      return "dev"
    } else if(role == 2) {
      return "guard"
    }
    return ""
  }

}
