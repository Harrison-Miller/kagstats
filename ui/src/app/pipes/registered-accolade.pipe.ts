import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'registeredAccolade',
  pure: true
})
export class RegisteredAccoladePipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return this.getRegisteredAccolade(value);
    return null;
  }

  parseDate(date: string): number {
    //2011-11-09 05:20:22
    var b = date.split(/\D/);
    return Date.UTC(+b[0], +b[1]-1, +b[2], +b[3], +b[4], +b[5])
  }

  getRegisteredAccolade(registered: string): string {
    let l = this.parseDate(registered)
    let n = Date.now();

    let diff = n - l;
    diff = diff / (1000 * 3600 * 24);

    if(diff > 10*365) {
      return "r10yBadge";
    } else if(diff > 9*365) {
      return "r9yBadge";
    } else if(diff > 8*365) {
      return "r8yBadge";
    } else if(diff > 7*365) {
      return "r7yBadge";
    } else if(diff > 6*365) {
      return "r6yBadge";
    } else if(diff > 5*365) {
      return "r5yBadge";
    } else if(diff > 4*365) {
      return "r4yBadge";
    } else if(diff > 3*365) {
      return"r3yBadge";
    } else if(diff > 2*365) {
      return "r2yBadge";
    } else if(diff > 1*365) {
      return "r1yBadge";
    } else if(diff > 9*31) {
      return "r9mBadge";
    } else if(diff > 6*31) {
      return "r6mBadge";
    } else if(diff > 3*31) {
      return "r3mBadge";
    } else if(diff > 2*31) {
      return "r2mBadge";
    } else if(diff > 1*31) {
      return "r1mBadge";
    } else if(diff > 3*7) {
      return "r3wBadge";
    } else if(diff > 2*7) {
      return "r2wBadge";
    } else if(diff > 1*7) {
      return "r1wBadge";
    }

    return "newBadge";
  }

}
