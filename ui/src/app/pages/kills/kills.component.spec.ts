import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KillsComponent } from './kills.component';

describe('KillsComponent', () => {
  let component: KillsComponent;
  let fixture: ComponentFixture<KillsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KillsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KillsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
