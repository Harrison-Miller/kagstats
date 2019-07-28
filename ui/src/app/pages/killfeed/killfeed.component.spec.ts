import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KillfeedComponent } from './killfeed.component';

describe('KillfeedComponent', () => {
  let component: KillfeedComponent;
  let fixture: ComponentFixture<KillfeedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KillfeedComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KillfeedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
