import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'statusClass',
  pure: true
})
export class StatusClassPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getStatusClass(value);
  }

  getStatusClass(lastUpdate: string): string {
    let l = Date.parse(lastUpdate + " UTC");
    let n = Date.now();

    let diff = n - l;

    if(diff > 600000) {
      return "bg-secondary"
    }
    if(diff > 300000){
      return "bg-warning"
    }
    return "bg-success"
  }

}
