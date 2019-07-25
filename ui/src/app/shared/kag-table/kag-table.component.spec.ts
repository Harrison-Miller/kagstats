import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KagTableComponent } from './kag-table.component';

describe('KagTableComponent', () => {
  let component: KagTableComponent;
  let fixture: ComponentFixture<KagTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KagTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KagTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
